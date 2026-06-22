package gitlab_dash_client

import (
	"context"
	"fmt"
	"gitlab_api/pkg/http_client"
	"gitlab_api/pkg/http_request_proxy"
	"slices"
)

type GitlabDashClient struct {
	client                      *http_client.HttpClient
	projectIDList               map[int]struct{}
	ignoreTestBranchCompareList map[int]struct{}
	testBranchName              string
}

func NewGitlabDashClient(
	client *http_client.HttpClient,
	projectIDList map[int]struct{},
	ignoreTestBranchCompareList map[int]struct{},
	testBranchName string,
) *GitlabDashClient {
	return &GitlabDashClient{
		client:                      client,
		projectIDList:               projectIDList,
		ignoreTestBranchCompareList: ignoreTestBranchCompareList,
		testBranchName:              testBranchName,
	}
}

func (c *GitlabDashClient) FindProjectsByIDList(ctx context.Context, idList []int) ([]Project, error) {
	data, err := http_request_proxy.HandleHttpClientRequest(
		func() ([]Project, error) {
			result := make([]Project, 0, len(c.projectIDList))

			for id := range c.projectIDList {
				var (
					data Project
					//
					path = fmt.Sprintf(project_path, id)
				)

				_, err := c.client.
					Get(ctx, path).
					SetDestination(&data).
					Execute()
				if err != nil {
					return result, err
				}

				result = append(result, data)
			}

			return result, nil
		},
	)
	if err != nil {
		return data, err
	}

	slices.SortFunc(data, func(a, b Project) int {
		return b.LastActivityDate.Compare(a.LastActivityDate)
	})

	return data, err
}

func (c *GitlabDashClient) FindUserInfo(ctx context.Context) (UserInfo, error) {
	return http_request_proxy.HandleHttpClientRequest(
		func() (UserInfo, error) {
			var data UserInfo

			_, err := c.client.
				Get(ctx, user_path).
				SetDestination(&data).
				Execute()

			return data, err
		},
	)
}

func (c *GitlabDashClient) FindProjectInfoMapByList(ctx context.Context, projectList []Project) (map[int]ProjectInfo, error) {
	result := make(map[int]ProjectInfo, len(projectList))

	for _, p := range projectList {
		data, err := c.findProjectInfoByID(ctx, p.ID, p.DefaultBranch)
		if err != nil {
			return nil, err
		}

		result[p.ID] = data
	}

	return result, nil
}

func (c *GitlabDashClient) findProjectInfoByID(ctx context.Context, id int, defaultBranchName string) (ProjectInfo, error) {
	var (
		zero                  ProjectInfo
		tagDisplayInfo        *BranchDisplayInfo
		testBranchDisplayData *BranchDisplayInfo
		defaultBranchIsActual *bool
	)

	tagData, err := http_request_proxy.HandleHttpClientRequest(
		func() (*BranchCommonData, error) {
			data := make([]BranchCommonData, 0, 1)

			_, err := c.client.
				Get(ctx, fmt.Sprintf(tag_path, id)).
				SetQuery(map[string]any{
					"per_page": 1,
				}).
				SetDestination(&data).
				Execute()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, nil
			}

			return &data[0], nil
		},
	)
	if err != nil {
		return zero, err
	}

	defaultBranchData, err := http_request_proxy.HandleHttpClientRequest(
		func() (BranchCommonData, error) {
			var data BranchCommonData

			_, err := c.client.
				Get(ctx, fmt.Sprintf(branch_path, id, defaultBranchName)).
				SetDestination(&data).
				Execute()
			return data, err
		},
	)
	if err != nil {
		return zero, err
	}

	_, exists := c.ignoreTestBranchCompareList[id]
	if !exists {
		testBranchData, err := http_request_proxy.HandleHttpClientRequest(
			func() (BranchCommonData, error) {
				var data BranchCommonData

				_, err := c.client.
					Get(ctx, fmt.Sprintf(branch_path, id, c.testBranchName)).
					SetDestination(&data).
					Execute()
				return data, err
			},
		)
		if err != nil {
			return zero, err
		}

		compareInfo, err := http_request_proxy.HandleHttpClientRequest(
			func() (BranchCompareData, error) {
				var data BranchCompareData

				_, err := c.client.
					Get(ctx, fmt.Sprintf(brach_compare_path, id)).
					SetQuery(map[string]any{
						"from": defaultBranchName,
						"to":   c.testBranchName,
					}).
					SetDestination(&data).
					Execute()
				return data, err
			},
		)
		if err != nil {
			return zero, err
		}

		testBranchDisplayData = &BranchDisplayInfo{
			CommitID:  testBranchData.Commit.ID,
			Name:      testBranchData.Name,
			UpdatedAt: testBranchData.Commit.UpdatedAt,
		}

		isActual := len(compareInfo.Diffs) == 2
		defaultBranchIsActual = &isActual
	}

	if tagData != nil {
		isAc := tagData.Commit.ID == defaultBranchData.Commit.ID
		tagDisplayInfo = &BranchDisplayInfo{
			CommitID:  tagData.Commit.ID,
			Name:      tagData.Name,
			UpdatedAt: tagData.Commit.UpdatedAt,
			IsActual:  &isAc,
		}
	}

	return ProjectInfo{
		TagInfo: tagDisplayInfo,
		DefaultBranchInfo: BranchDisplayInfo{
			CommitID:  defaultBranchData.Commit.ID,
			Name:      defaultBranchData.Name,
			UpdatedAt: defaultBranchData.Commit.UpdatedAt,
			IsActual:  defaultBranchIsActual,
		},
		TestBranchInfo: testBranchDisplayData,
	}, nil
}
