package model

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
)

type ADOConnectionInfo struct {
	ConnectionUrl string `json:"connection_url,omitempty"`
	ConnectionPAT string `json:"connection_pat,omitempty"`
}

type GitRepo struct {
	Name          string      `json:"name"`
	Id            uuid.UUID   `json:"id"`
	Url           string      `json:"url"`
	WebUrl        string      `json:"webUrl"`
	Links         interface{} `json:"links"`
	DefaultBranch string      `json:"defaultBranch"`
}

type GitCommitsCriteria struct {
	RepositoryId string `json:"repository_id,omitempty"`
	Author       string `json:"author,omitempty"`
	User         string `json:"user,omitempty"`
	FromDate     string `json:"from_date,omitempty"`
	Version      string `json:"version,omitempty"`
	VersionType  string `json:"version_type,omitempty"`
	Skip         int    `json:"skip,omitempty"`
	Top          int    `json:"stop,omitempty"`
}

type GitCommitItem struct {
	Repository string             `json:"repository,omitempty"`
	CommitInfo []git.GitCommitRef `json:"commit_info,omitempty"`
}

type GitCommitItemSimple struct {
	Repository string             `json:"repository,omitempty"`
	CommitInfo []CommitInfoSimple `json:"commit_info,omitempty"`
}

type CommitInfoSimple struct {
	Author           git.GitUserDate      `json:"author,omitempty"`
	Comment          string               `json:"comment,omitempty"`
	CommentTruncated bool                 `json:"comment_truncated,omitempty"`
	CommitId         string               `json:"commit_id,omitempty"`
	Committer        git.GitUserDate      `json:"committer,omitempty"`
	Changes          []interface{}        `json:"changes,omitempty"`
	Push             git.GitPushRef       `json:"push,omitempty"`
	RemoteUrl        string               `json:"remote_url,omitempty"`
	WorkItems        []webapi.ResourceRef `json:"work_items,omitempty"`
}
