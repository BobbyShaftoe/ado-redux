package model

import (
	"HTTP_Sever/helpers/ado"
)

type HelloData struct {
	Name string `json:"name"`
}

type DashboardData struct {
	Projects []string        `json:"projects"`
	Repos    []ado.GitRepo   `json:"repos"`
	Commits  []GitCommitItem `json:"commits"`
}
