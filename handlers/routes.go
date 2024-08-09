package handlers

import (
	"HTTP_Sever/helpers/ado"
	"HTTP_Sever/views"
	"github.com/a-h/templ"
)

type adoRouteInterface interface {
	RenderRouteTempl(templMap TemplMap, route string) templ.Component
}

type GitRepo []ado.GitRepo

type TemplMap struct {
	templMap map[string]templ.Component
}

func NewTemplMap() *TemplMap {
	return &TemplMap{
		templMap: map[string]templ.Component{
			"index":     views.Layout(),
			"hello":     views.Hello(),
			"dashboard": views.Dashboard([]ado.GitRepo{}),
		},
	}
}

func (g GitRepo) RenderRouteTempl(templMap TemplMap, route string, gitRepo GitRepo) templ.Component {
	return templMap.templMap[route]
}

func RenderRouteTempl(templMap TemplMap, route string) templ.Component {
	return templMap.templMap[route]
}
