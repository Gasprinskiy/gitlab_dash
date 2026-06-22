package main

import (
	"context"
	"gitlab_api/config"
	"gitlab_api/core/gitlab_dash_client"
	"gitlab_api/pkg/http_client"
	"gitlab_api/pkg/s_call"
	"gitlab_api/ui_components"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ReloadMsg struct{}

type Model struct {
	projectList       []gitlab_dash_client.Project
	projectInfoMap    map[int]gitlab_dash_client.ProjectInfo
	cursorProjectInfo gitlab_dash_client.ProjectInfo
	//
	cfg        *config.Config
	userInfo   *gitlab_dash_client.UserInfo
	dashClient *gitlab_dash_client.GitlabDashClient
	table      table.Model

	width  int
	height int
}

func initialModel() Model {
	cfg := config.NewConfig()

	projectIDMap := make(map[int]struct{}, len(cfg.ProjectsData.ProjectIdList))
	ignoreTestBranchCompareMap := make(map[int]struct{}, len(cfg.ProjectsData.IgnoreTestBranchCompareList))

	for _, id := range cfg.ProjectsData.ProjectIdList {
		projectIDMap[id] = struct{}{}
	}

	for _, id := range cfg.ProjectsData.IgnoreTestBranchCompareList {
		ignoreTestBranchCompareMap[id] = struct{}{}
	}

	m := Model{
		cfg:   cfg,
		table: table.New(table.WithFocused(true), table.WithStyles(ui_components.TableStyles())),
		dashClient: gitlab_dash_client.NewGitlabDashClient(
			http_client.NewHttpClient(
				cfg.Credentials.Host,
				cfg.Credentials.PersonalToken,
			),
			projectIDMap,
			ignoreTestBranchCompareMap,
			cfg.ProjectsData.TestBranchName,
		),
	}

	m.fetchProjectInfo()

	return m
}

func (m *Model) fetchProjectInfo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	userInfo, err := m.dashClient.FindUserInfo(ctx)
	if err != nil {
		log.Fatal("could not find user info: ", err)
	}

	projectsList, err := m.dashClient.FindProjectsByIDList(ctx, nil)
	if err != nil {
		log.Fatal("could not find projects list: ", err)
	}

	projectInfoMap, err := m.dashClient.FindProjectInfoMapByList(ctx, projectsList)
	if err != nil {
		log.Fatal("could not find projects info map: ", err)
	}

	m.userInfo = &userInfo
	m.projectList = projectsList
	m.projectInfoMap = projectInfoMap
}

func (m *Model) setTableRows() {
	if len(m.table.Columns()) == 0 {
		return
	}

	rows := make([]table.Row, 0, len(m.projectList))
	for _, p := range m.projectList {
		rows = append(rows, table.Row{
			strconv.Itoa(p.ID),
			p.Name,
			p.DefaultBranch,
			p.LastActivityDate.In(time.Local).Format(time.DateTime),
		})
	}
	m.table.SetRows(rows)
}

func (m *Model) setCurrentCursorProjectInfo() {
	idx := m.table.Cursor()
	id := m.projectList[idx].ID

	m.cursorProjectInfo = m.projectInfoMap[id]
}

func (m Model) userInfoBox() string {
	key := lipgloss.NewStyle().Foreground(lipgloss.Color(ui_components.COLOR_TEXT_DEFAULT)).Bold(true)
	val := lipgloss.NewStyle().Foreground(lipgloss.Color(ui_components.COLOR_TEXT_PRIMARY))

	u := m.userInfo

	rows := []string{
		key.Render("Host:            ") + val.Render(m.cfg.Credentials.Host),
		key.Render(""),
		key.Render("UserName:        ") + val.Render(u.UserName),
		key.Render(""),
		key.Render("Last activity:   ") + val.Render(u.LastActivityOn),
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) branchInfoBox(info *gitlab_dash_client.BranchDisplayInfo) string {
	key := lipgloss.NewStyle().Foreground(lipgloss.Color(ui_components.COLOR_TEXT_DEFAULT)).Bold(true)
	val := lipgloss.NewStyle().Foreground(lipgloss.Color(ui_components.COLOR_TEXT_PRIMARY))
	if info == nil {
		rows := []string{
			key.Render("Not found"),
		}

		return lipgloss.JoinVertical(lipgloss.Center, rows...)
	}

	rows := []string{
		key.Render("Name:          ") + val.Render(info.Name),
		key.Render(""),
		key.Render("Fresh:         ") + val.Render(info.IsActual),
		key.Render(""),
		key.Render("Commit date:   ") + val.Render(info.UpdatedAt.In(time.Local).Format(time.DateTime)),
		key.Render(""),
		key.Render("Commit ID:     ") + val.Render(info.CommitID),
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *Model) tableLayout() {

	idW := m.width * 5 / 100
	titleW := m.width * 20 / 100
	branchW := m.width * 10 / 100
	activityW := m.width * 15 / 100

	m.table.SetColumns([]table.Column{
		{Title: "ID", Width: idW},
		{Title: "Name", Width: titleW},
		{Title: "Default branch", Width: branchW},
		{Title: "Last Activity", Width: activityW},
	})
	m.setTableRows()
	m.setCurrentCursorProjectInfo()

	bottomH := m.height - m.height*25/100
	tableH := bottomH - 4
	if tableH < 1 {
		tableH = 1
	}
	m.table.SetHeight(tableH)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.tableLayout()

	case tea.KeyMsg:
		msgStr := msg.String()

		if msgStr == "q" || msgStr == "ctrl+c" {
			return m, tea.Quit
		}

		if msgStr == "ctrl+r" {
			m.fetchProjectInfo()

			return m, func() tea.Msg {
				return ReloadMsg{}
			}
		}

		if msgStr == "enter" {
			idx := m.table.Cursor()
			if idx >= 0 && idx <= (len(m.projectList)-1) {
				s_call.OpenBrowser(m.projectList[idx].WebUrl)
			}
		}

	case ReloadMsg:
		m.setTableRows()
	}

	prev := m.table.Cursor() // где стояли ДО

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg) // навигация сдвинула курсор

	cmds := []tea.Cmd{cmd}

	if m.table.Cursor() != prev {
		m.setCurrentCursorProjectInfo()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	topH := m.height * 22 / 100
	bottomH := m.height - topH
	rightWidth := m.width * 43 / 100

	top := ui_components.Box("Common info", m.userInfoBox(), m.width*57/100, topH, ui_components.TitleCenter)
	bottom := ui_components.Box("Repositories", m.table.View(), m.width*57/100, bottomH, ui_components.TitleCenter)
	right := ui_components.Box(
		"Repository info",
		lipgloss.JoinVertical(
			lipgloss.Top,
			ui_components.Box(
				"Last tag",
				m.branchInfoBox(m.cursorProjectInfo.TagInfo),
				rightWidth-4,
				m.height*20/100,
				ui_components.TitleLeft,
			),
			ui_components.Box(
				"Default branch",
				m.branchInfoBox(&m.cursorProjectInfo.DefaultBranchInfo),
				rightWidth-4,
				m.height*20/100,
				ui_components.TitleLeft,
			),
			ui_components.Box(
				"Test branch",
				m.branchInfoBox(m.cursorProjectInfo.TestBranchInfo),
				rightWidth-4,
				m.height*20/100,
				ui_components.TitleLeft,
			),
		),
		rightWidth,
		m.height,
		ui_components.TitleCenter,
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, top, bottom), right)
}

func init() {
	runtime.GOMAXPROCS(8)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
