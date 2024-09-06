package model

import (
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
)

type HelloData struct {
	Name string `json:"name"`
}

type DashboardData struct {
	Projects []string        `json:"projects"`
	Repos    []GitRepo       `json:"repos"`
	Commits  []GitCommitItem `json:"commits"`
	Builds   []build.Build   `json:"builds"`
}
