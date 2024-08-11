package ado

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"log"
)

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
		log.Printf("Project Name[%v] = %v", index, *teamProjectReference.Name)
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
		log.Printf("Repository Name[%v] = %v", index, *gitRepository.Name)
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
