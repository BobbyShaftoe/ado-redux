package model

import "HTTP_Sever/helpers/ado"

type RepositoriesData struct {
	Projects []string      `json:"projects"`
	Repos    []ado.GitRepo `json:"repos"`
}
