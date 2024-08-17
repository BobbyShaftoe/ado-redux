package handlers

import (
	"HTTP_Sever/helpers/ado"
	"HTTP_Sever/model"
	"HTTP_Sever/views"
	"context"
	"github.com/a-h/templ"
	"net/http"
)

func RenderHello(globalState *model.GlobalState) templ.Component {
	logger.json.Info("RenderHello", "globalState", globalState)
	return views.Hello(globalState)
}

func RenderIndex(globalState *model.GlobalState) templ.Component {
	logger.json.Info("RenderIndex", "globalState", globalState)
	return views.Index(globalState)
}

func RenderDashboard(globalState *model.GlobalState) templ.Component {
	logger.json.Info("RenderDashboard", "globalState", globalState)

	dashboardData := getDashboardData(globalState)
	globalState.UpdateGlobalStateProjects(dashboardData.Projects)
	return views.Dashboard(dashboardData, globalState)
}

func RenderDashboardUpdateProject(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		logger.json.Info("RenderDashboardUpdateProject", "globalState", globalState)

		globalState.UpdateGlobalStateProject(project)
		dashboardData := getDashboardData(globalState)
		templ.Handler(views.DashboardContent(dashboardData, globalState)).ServeHTTP(w, r)
	}
}

func getDashboardData(globalState *model.GlobalState) model.DashboardData {
	adoCtx := context.Background()

	projects := NewADOClients(adoCtx).GetProjects(adoCtx)
	repositories := NewADOClients(adoCtx).GetRepositories(adoCtx, globalState.CurrentProject)

	dashboardData := model.DashboardData{
		Projects: ado.ReturnProjects(projects),
		Repos:    ado.ReturnGitRepos(repositories),
	}
	logger.json.Info("getDashboardData", "dashboardData", dashboardData, "globalState", globalState)
	return dashboardData
}
