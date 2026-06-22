package ui_components

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	COLOR_BORDER        = "#4dad8d"
	COLOR_TEXT_DEFAULT  = "#ffffffe6"
	COLOR_TEXT_PRIMARY  = "#63e2b7"
	SELECTED_BACKGROUND = "#233633"
)

type TitlePos int

const (
	TitleLeft TitlePos = iota
	TitleCenter
	TitleRight
)

func Box(title, content string, w, h int, pos TitlePos) string {
	bc := lipgloss.NewStyle().Foreground(lipgloss.Color(COLOR_BORDER)) //
	b := lipgloss.RoundedBorder()

	style := lipgloss.NewStyle().
		Border(b).
		BorderForeground(lipgloss.Color(COLOR_BORDER)).
		Width(w - 2).Height(h - 2).
		Padding(1)

	inner := w - 2
	if inner < 0 {
		inner = 0
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color(COLOR_TEXT_DEFAULT)).
		Render(" " + title + " ")
	labelW := lipgloss.Width(label)

	if labelW > inner {
		label, labelW = "", 0
	}

	remaining := inner - labelW

	var leftPad, rightPad int
	switch pos {
	case TitleCenter:
		leftPad = remaining / 2
		rightPad = remaining - leftPad
	case TitleRight:
		rightPad = 1
		if rightPad > remaining {
			rightPad = remaining
		}
		leftPad = remaining - rightPad
	default:
		leftPad = 1
		if leftPad > remaining {
			leftPad = remaining
		}
		rightPad = remaining - leftPad
	}

	top := bc.Render(b.TopLeft) +
		bc.Render(strings.Repeat(b.Top, leftPad)) +
		label +
		bc.Render(strings.Repeat(b.Top, rightPad)) +
		bc.Render(b.TopRight)

	lines := strings.Split(style.Render(content), "\n")
	lines[0] = top
	return strings.Join(lines, "\n")
}

func TableStyles() table.Styles {
	s := table.DefaultStyles()

	s.Header = s.Header.
		Foreground(lipgloss.Color(COLOR_TEXT_DEFAULT)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(COLOR_BORDER)).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color(COLOR_TEXT_DEFAULT)).
		Background(lipgloss.Color(SELECTED_BACKGROUND)).
		Bold(true)

	return s
}
