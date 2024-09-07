package ado

import (
	"HTTP_Sever/model"
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"log/slog"
	"os"
	"regexp"
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

func ReturnProjects(responseValue *core.GetProjectsResponseValue) []string {
	var Projects []string
	index := 0
	// Log the page of team project names
	for _, teamProjectReference := range (*responseValue).Value {
		logger.json.Debug("ReturnProjects", "index", index, "projectName", *teamProjectReference.Name)
		Projects = append(Projects, *teamProjectReference.Name)
		index++
	}
	return Projects
}

func ReturnGitRepos(responseValue *[]git.GitRepository) []model.GitRepo {
	var Repositories []model.GitRepo
	index := 0
	// Log the page of team project names
	for _, gitRepository := range *responseValue {
		logger.json.Debug("ReturnGitRepos", "index", index, "gitRepository", *gitRepository.Name)
		Repositories = append(Repositories, model.GitRepo{
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

func ReturnGitRepoNames(gitRepositories *[]git.GitRepository) []string {
	var repositories []string
	index := 0
	for _, gitRepository := range *gitRepositories {
		logger.json.Debug("ReturnGitRepos", "index", index, "gitRepository", *gitRepository.Name)
		repositories = append(repositories, *gitRepository.Name)
		index++
	}
	return repositories
}

func ReturnGitCommitItemSimple(gitCommitItem []model.GitCommitItem) []model.GitCommitItemSimple {
	logger.json.Info("ReturnGitCommitItemSimple", "gitCommitItem", gitCommitItem)
	commitItems := make([]model.GitCommitItemSimple, 0)
	ci := make([]model.CommitInfoSimple, 0)

	for _, commitItem := range gitCommitItem {
		ci = ci[:0]

		for _, commitInfo := range commitItem.CommitInfo {
			tmpCommitInfoSimple := model.CommitInfoSimple{
				Author:    *commitInfo.Author,
				Comment:   *commitInfo.Comment,
				CommitId:  *commitInfo.CommitId,
				Committer: *commitInfo.Committer,
				Push:      *commitInfo.Push,
				RemoteUrl: *commitInfo.RemoteUrl,
			}

			if commitInfo.CommentTruncated != nil {
				tmpCommitInfoSimple.CommentTruncated = *commitInfo.CommentTruncated
			}
			if commitInfo.Changes != nil {
				tmpCommitInfoSimple.Changes = *commitInfo.Changes
			}
			if commitInfo.WorkItems != nil {
				tmpCommitInfoSimple.WorkItems = *commitInfo.WorkItems
			}

			tmpCommitInfoSimple.CommitIdShort = (*commitInfo.CommitId)[:7]

			repoUrlPattern := regexp.MustCompile("(https://.*?)/commit/[a-f0-9]+$")
			tmpCommitInfoSimple.RemoteRepoUrl = repoUrlPattern.ReplaceAllString(*commitInfo.RemoteUrl, "$1")

			logger.json.Info("R", "url", tmpCommitInfoSimple.RemoteRepoUrl, "comment", tmpCommitInfoSimple.Comment)
			ci = append(ci, tmpCommitInfoSimple)
		}
		commitItems = append(commitItems, model.GitCommitItemSimple{
			Repository: commitItem.Repository,
			CommitInfo: ci,
		})
	}
	logger.json.Debug("ReturnGitCommitItemSimple", "commitItems", commitItems)
	return commitItems
}

type UserValidators interface {
	ValidateUser(user string, userGraph *[]graph.GraphUser) bool
	ValidateUsers(users []string, userGraph *[]graph.GraphUser) bool
}

type UserValidator struct {
	SystemUsers *[]graph.GraphUser
}

func NewUserValidator(systemUsers *[]graph.GraphUser) *UserValidator {
	return &UserValidator{SystemUsers: systemUsers}
}

func (uv *UserValidator) ValidateUsers(users []string) (bool, []string) {
	var validUsers = make([]string, 0)
	logger.json.Info("ValidateUsers", "lengthUsers", len(users), "users", users)
	if len(users) == 0 {
		return false, validUsers
	}
	for _, user := range users {
		res, user := uv.ValidateUser(user)
		if !res {
			logger.json.Info("ValidateUser", "invalidUser", user)
			continue
		}
		validUsers = append(validUsers, user)
	}
	return true, validUsers
}

func (uv *UserValidator) ValidateUser(user string) (bool, string) {
	logger.json.Debug("ValidateUser", "users", uv.SystemUsers)
	for _, graphUser := range *uv.SystemUsers {
		logger.json.Info("ValidateUser", "user", user, "graphUser", *graphUser.MailAddress)
		if *graphUser.PrincipalName == user || *graphUser.MailAddress == user {
			return true, user
		}
	}
	return false, user
}
