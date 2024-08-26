package model

import "github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"

type ADOConnectionInfo struct {
	ConnectionUrl string `json:"connection_url,omitempty"`
	ConnectionPAT string `json:"connection_pat,omitempty"`
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
