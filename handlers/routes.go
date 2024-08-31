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
	//logger.json.Info("RenderHello", "globalState", globalState)
	return views.Hello(globalState)
}

func RenderIndex(globalState *model.GlobalState) templ.Component {
	//logger.json.Info("RenderIndex", "globalState", globalState)
	return views.Index(globalState)
}

func RenderDashboardHandler(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.json.Debug("RenderDashboard", "globalState", globalState)
		logger.json.Info("RenderDashboard", "project", globalState.CurrentProject)

		dashboardData := getDashboardData(globalState)
		globalState.UpdateGlobalStateProjects(dashboardData.Projects)
		templ.Handler(views.Dashboard(dashboardData, globalState)).ServeHTTP(w, r)
	}
}

func RenderRepositoriesHandler(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.json.Debug("RenderDashboard", "globalState", globalState)
		logger.json.Info("RenderRepositories", "project", globalState.CurrentProject)

		repositoriesData := getRepositoriesData(globalState)
		globalState.UpdateGlobalStateProjects(repositoriesData.Projects)
		templ.Handler(views.Repositories(repositoriesData, globalState)).ServeHTTP(w, r)
		//return views.Repositories(repositoriesData, globalState)
	}
}

func RenderDashboardUpdateProject(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		logger.json.Debug("RenderDashboardUpdateProject", "globalState", globalState)
		//logger.json.Info("RenderDashboardUpdateProject", "project", globalState.CurrentProject)

		globalState.UpdateGlobalStateProject(project)
		dashboardData := getDashboardData(globalState)
		templ.Handler(views.DashboardContent(dashboardData, globalState)).ServeHTTP(w, r)
	}
}

func RenderRepositoriesUpdateProject(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		logger.json.Debug("RenderRepositoriesUpdateProject", "globalState", globalState)
		logger.json.Info("RenderDashboardUpdateProject", "project", globalState.CurrentProject)

		globalState.UpdateGlobalStateProject(project)
		logger.json.Info("RenderDashboardUpdateProject", "project", globalState.CurrentProject)
		repositoriesData := getRepositoriesData(globalState)
		templ.Handler(views.RepositoriesContent(repositoriesData, globalState)).ServeHTTP(w, r)
	}
}

func getDashboardData(globalState *model.GlobalState) model.DashboardData {
	adoCtx := context.Background()

	projects := NewADOClients(adoCtx).GetProjects(adoCtx)
	repoNames := NewADOClients(adoCtx).GetRepositories(adoCtx, globalState)
	repositories := ado.ReturnGitRepoNames(repoNames)
	commitsCriteria := ReturnGitCommitCriteria(globalState)
	commits := NewADOClients(adoCtx).GetCommits(adoCtx, commitsCriteria, globalState, repositories)

	//logger.json.Info("getDashboardData", "commitsCriteria", commitsCriteria, "commits", commits)

	dashboardData := model.DashboardData{
		Projects: ado.ReturnProjects(projects),
		Repos:    ado.ReturnGitRepos(repoNames),
		Commits:  commits,
	}
	//logger.json.Info("getDashboardData", "dashboardData", dashboardData, "globalState", globalState)
	return dashboardData
}

func getRepositoriesData(globalState *model.GlobalState) model.RepositoriesData {
	adoCtx := context.Background()

	projects := NewADOClients(adoCtx).GetProjects(adoCtx)
	repoNames := NewADOClients(adoCtx).GetRepositories(adoCtx, globalState)

	//logger.json.Info("getDashboardData", "commitsCriteria", commitsCriteria, "commits", commits)

	repositoriesData := model.RepositoriesData{
		Projects: ado.ReturnProjects(projects),
		Repos:    ado.ReturnGitRepos(repoNames),
	}
	//logger.json.Info("getDashboardData", "dashboardData", dashboardData, "globalState", globalState)
	return repositoriesData
}
