package ado

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"log/slog"
	"os"
)

type localLogger struct {
	json *slog.Logger
}

var logger = &localLogger{
	json: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
}

func fatal(v ...any) {
	logger.json.Error("main", "err", fmt.Sprint(v...))
	os.Exit(1)
}

type GitRepo struct {
	Name          string      `json:"name"`
	Id            uuid.UUID   `json:"id"`
	Url           string      `json:"url"`
	WebUrl        string      `json:"webUrl"`
	Links         interface{} `json:"links"`
	DefaultBranch string      `json:"defaultBranch"`
}

func ReturnProjects(responseValue *core.GetProjectsResponseValue) []string {
	var Projects []string
	index := 0
	// Log the page of team project names
	for _, teamProjectReference := range (*responseValue).Value {
		logger.json.Info("ReturnProjects", "index", index, "projectName", *teamProjectReference.Name)
		Projects = append(Projects, *teamProjectReference.Name)
		index++
	}
	return Projects
}

func ReturnGitRepos(responseValue *[]git.GitRepository) []GitRepo {
	var Repositories []GitRepo
	index := 0
	// Log the page of team project names
	for _, gitRepository := range *responseValue {
		logger.json.Info("ReturnGitRepos", "index", index, "gitRepository", *gitRepository.Name)
		Repositories = append(Repositories, GitRepo{
			Name:          *gitRepository.Name,
			Id:            *gitRepository.Id,
			Url:           *gitRepository.Url,
			WebUrl:        *gitRepository.WebUrl,
			Links:         gitRepository.Links,
			DefaultBranch: *gitRepository.DefaultBranch,
		})
		index++
	}
	return Repositories
}

func validateUser(user string, userGraph *[]graph.GraphUser) bool {
	for _, graphUser := range *userGraph {
		if *graphUser.MailAddress == user {
		}
	}
	return false
}
