package handlers

import (
	"HTTP_Sever/helpers/ado"
	"HTTP_Sever/model"
	"HTTP_Sever/views"
	"context"
	"encoding/json"
	"github.com/a-h/templ"
	"net/http"
	"regexp"
)

func RenderIndex(globalState *model.GlobalState) templ.Component {
	logger.json.Debug("RenderIndex", "globalState", globalState)
	return views.Index(globalState)
}

func HandleSearch(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.json.Info("HandleSearch", "search", r.FormValue("search"))

		currentPage := r.FormValue("pageContext")
		pageData := r.FormValue("modelData")
		searchTerm := r.FormValue("search")

		if currentPage == "repositories" {
			logger.json.Info("HandleSearch", "match", "repositories", "page", currentPage)
			repositoriesData := filterRepositories(pageData, searchTerm)
			templ.Handler(views.RepositoriesContent(repositoriesData, globalState)).ServeHTTP(w, r)
		}

		if currentPage == "dashboard" {
			logger.json.Info("HandleSearch", "match", "dashboard", "page", currentPage)
			dashboardData := filterDashboard(pageData, searchTerm)
			templ.Handler(views.DashboardContents(dashboardData)).ServeHTTP(w, r)

		}
	}
}

func filterDashboard(dashboardData string, search string) model.DashboardData {
	dashboardStruct := model.DashboardData{}
	dashboardBytes := []byte(dashboardData)
	_ = json.Unmarshal(dashboardBytes, &dashboardStruct)

	filteredDashboardCommitItems := make([]model.GitCommitItemSimple, 0)
	matchPattern, err := regexp.Compile(".*" + search + ".*")

	for _, commit := range dashboardStruct.Commits {

		if err != nil {
			logger.json.Error("filterDashboard", "err", err)
			break
		}

		matched := false

		if m, err := regexp.Match(matchPattern.String(), []byte(commit.CommitInfo[0].Comment)); err == nil && m {
			matched = true
		}
		if m, err := regexp.Match(matchPattern.String(), []byte(commit.Repository)); err == nil && m {
			matched = true
		}

		if matched {
			filteredDashboardCommitItems = append(filteredDashboardCommitItems, commit)
		}

	}
	dashboardStruct.Commits = filteredDashboardCommitItems
	return dashboardStruct
}

func filterRepositories(repositoriesData string, search string) model.RepositoriesData {
	repositoriesStruct := model.RepositoriesData{}
	repositoriesBytes := []byte(repositoriesData)
	_ = json.Unmarshal(repositoriesBytes, &repositoriesStruct)

	filteredRepos := make([]model.GitRepo, 0)
	matchPattern, err := regexp.Compile(".*" + search + ".*")

	for _, repo := range repositoriesStruct.Repos {
		if err != nil {
			logger.json.Error("filterRepositories", "err", err)
			break
		}
		if matched, err := regexp.Match(matchPattern.String(), []byte(repo.Name)); err == nil && matched {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	repositoriesStruct.Repos = filteredRepos
	return repositoriesStruct
}

func RenderDashboardHandler(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.json.Debug("RenderDashboard", "globalState", globalState)

		dashboardData := getDashboardData(globalState)
		globalState.UpdateGlobalStateProjects(dashboardData.Projects)
		templ.Handler(views.Dashboard(dashboardData, globalState)).ServeHTTP(w, r)
	}
}

func RenderRepositoriesHandler(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.json.Debug("RenderDashboard", "globalState", globalState)

		repositoriesData := getRepositoriesData(globalState)
		globalState.UpdateGlobalStateProjects(repositoriesData.Projects)
		templ.Handler(views.Repositories(repositoriesData, globalState)).ServeHTTP(w, r)
	}
}

func RenderDashboardUpdateProject(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		logger.json.Debug("RenderDashboardUpdateProject", "globalState", globalState)

		globalState.UpdateGlobalStateProject(project)
		dashboardData := getDashboardData(globalState)
		templ.Handler(views.DashboardMain(dashboardData, globalState)).ServeHTTP(w, r)
	}
}

func RenderRepositoriesUpdateProject(globalState *model.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		logger.json.Debug("RenderRepositoriesUpdateProject", "globalState", globalState)

		globalState.UpdateGlobalStateProject(project)
		repositoriesData := getRepositoriesData(globalState)
		templ.Handler(views.RepositoriesContent(repositoriesData, globalState)).ServeHTTP(w, r)
	}
}

func getDashboardData(globalState *model.GlobalState) model.DashboardData {
	adoCtx := context.Background()

	allUsers := append(globalState.AdditionalUsers, globalState.User)
	allCommits := make([]model.GitCommitItem, 0)

	projects := NewADOClients(adoCtx).GetProjects(adoCtx)
	// Get repositories for the current project
	repoNames := NewADOClients(adoCtx).GetRepositories(adoCtx, globalState)
	repositories := ado.ReturnGitRepoNames(repoNames)
	commitsCriteria := ReturnGitCommitCriteria(globalState)

	for _, user := range allUsers {
		commitsCriteria.Author = user
		commits := NewADOClients(adoCtx).GetCommits(adoCtx, commitsCriteria, globalState, repositories)
		allCommits = append(allCommits, commits...)
	}
	logger.json.Debug("getDashboardData", "allCommits", allCommits)
	allCommitsSimple := ado.ReturnGitCommitItemSimple(allCommits)
	logger.json.Debug("getDashboardData", "allCommitsSimple", allCommitsSimple)

	dashboardData := model.DashboardData{
		Projects: ado.ReturnProjects(projects),
		Repos:    ado.ReturnGitRepos(repoNames),
		Commits:  allCommitsSimple,
	}
	logger.json.Debug("getDashboardData", "dashboardData", dashboardData, "globalState", globalState)

	// Get builds for the current project
	buildList := NewADOClients(adoCtx).ListBuilds(adoCtx, globalState, ado.ReturnGitRepos(repoNames))
	//logger.json.Info("getDashboardData", "buildList", buildList)
	latestBuilds := NewADOClients(adoCtx).GetLatestBuildsFromBuilds(adoCtx, globalState, buildList)

	//latestBuilds := NewADOClients(adoCtx).GetLatestBuildsFromRepositories(adoCtx, globalState, repositories)
	logger.json.Info("getDashboardData", "latestBuilds", latestBuilds)
	dashboardData.Builds = latestBuilds

	return dashboardData
}

// Repositories functions
// getRepositoriesData returns the repositories data for the current project
func getRepositoriesData(globalState *model.GlobalState) model.RepositoriesData {
	adoCtx := context.Background()

	projects := NewADOClients(adoCtx).GetProjects(adoCtx)
	repoNames := NewADOClients(adoCtx).GetRepositories(adoCtx, globalState)

	logger.json.Debug("getRepositoriesData", "repoNames", repoNames)

	repositoriesData := model.RepositoriesData{
		Projects: ado.ReturnProjects(projects),
		Repos:    ado.ReturnGitRepos(repoNames),
	}
	logger.json.Debug("getDashboardData", "repositoriesData", repositoriesData)

	_ = NewADOClients(adoCtx).ListBuilds(adoCtx, globalState, ado.ReturnGitRepos(repoNames))

	return repositoriesData
}
