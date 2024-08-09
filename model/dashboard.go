package model

import (
	"HTTP_Sever/helpers/ado"
)

type HelloData struct {
	Name string
}

type DashboardData struct {
	Projects []string
	Repos    []ado.GitRepo
}
