package handlers

import (
	"HTTP_Sever/helpers/ado"
	"HTTP_Sever/model"
	"HTTP_Sever/views"
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
)

func RenderHello(globalState *model.GlobalState) templ.Component {
	gs, _ := json.MarshalIndent(globalState, "", "\t")
	fmt.Println("RenderHello")
	fmt.Println(string(gs))
	return views.Hello(globalState)
}

func RenderIndex(globalState *model.GlobalState) templ.Component {
	gs, _ := json.MarshalIndent(globalState, "", "\t")
	fmt.Println("RenderIndex")
	fmt.Println(string(gs))
	return views.Index(globalState)
}

func RenderDashboard(globalState *model.GlobalState) templ.Component {
	gs, _ := json.MarshalIndent(*globalState, "", "\t")
	fmt.Println("RenderDashboard")
	fmt.Println(string(gs))

	dashboardData := getDashboardData()
	globalState.UpdateGlobalStateProjects(dashboardData.Projects)
	return views.Dashboard(dashboardData, globalState)
}

func RenderDashboardUpdateProject(globalState *model.GlobalState) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		globalState.UpdateGlobalStateProject(project)
		gs, _ := json.MarshalIndent(globalState, "", "\t")
		fmt.Println("RenderDashboardUpdateProject")
		fmt.Println(string(gs))
		templ.Handler(RenderDashboard(&*globalState)).ServeHTTP(w, r)
	})
}

func getDashboardData() model.DashboardData {
	adoCtx := context.Background()
	adoClients := NewADOClients(adoCtx)
	projects := adoClients.GetProjects(adoCtx)
	repositories := adoClients.GetRepositories(adoCtx)
	dashboardData := model.DashboardData{
		Projects: ado.ReturnProjects(projects),
		Repos:    ado.ReturnGitRepos(repositories),
	}
	return dashboardData
}
