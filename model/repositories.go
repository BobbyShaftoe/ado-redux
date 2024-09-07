package model

type RepositoriesData struct {
	Projects []string  `json:"projects"`
	Repos    []GitRepo `json:"repos"`
}
