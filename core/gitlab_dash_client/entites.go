package gitlab_dash_client

import (
	"encoding/json"
	"time"
)

type Project struct {
	LastActivityDate time.Time `json:"last_activity_at"`
	Name             string    `json:"name"`
	WebUrl           string    `json:"web_url"`
	DefaultBranch    string    `json:"default_branch"`
	ID               int       `json:"id"`
}

type UserInfo struct {
	UserName       string `json:"username"`
	LastActivityOn string `json:"last_activity_on"`
	ID             int    `json:"id"`
}

type Commit struct {
	UpdatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

type BranchCommonData struct {
	Commit Commit `json:"commit"`
	Name   string `json:"name"`
}

type BranchCompareData struct {
	Diffs json.RawMessage `json:"diffs"`
}

type BranchDisplayInfo struct {
	UpdatedAt time.Time
	Name      string
	IsActual  string
	CommitID  string
}

type ProjectInfo struct {
	TagInfo           *BranchDisplayInfo
	TestBranchInfo    *BranchDisplayInfo
	DefaultBranchInfo BranchDisplayInfo
}
