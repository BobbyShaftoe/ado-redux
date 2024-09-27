package ado

import (
	"HTTP_Sever/model"
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
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

func ReturnLinks(i interface{}) map[string]string {
	links := make(map[string]string)
	for k, v := range i.(map[string]interface{}) {
		linkRef := v.(map[string]interface{})["href"]
		links[k] = linkRef.(string)
	}
	return links
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

func ReturnGitCommitItemSimple(gitCommitItems []model.GitCommitItem) []model.GitCommitItemSimple {
	commitItems := make([]model.GitCommitItemSimple, 0)
	repoUrlPattern := regexp.MustCompile("(https://.*?)/commit/[a-f0-9]+$")

	for _, commitItem := range gitCommitItems {
		ci := make([]model.CommitInfoSimple, 0)

		for _, commitInfo := range commitItem.CommitInfo {
			commitInfoAuthor := *commitInfo.Author
			commitInfoComment := *commitInfo.Comment
			commitInfoCommitId := *commitInfo.CommitId
			commitInfoCommitter := *commitInfo.Committer
			commitInfoPush := *commitInfo.Push
			commitInfoRemoteUrl := *commitInfo.RemoteUrl

			var commitInfoCommentTruncated bool
			if commitInfo.CommentTruncated != nil {
				commitInfoCommentTruncated = *commitInfo.CommentTruncated
			}
			var commitInfoChanges []interface{}
			if commitInfo.Changes != nil {
				commitInfoChanges = *commitInfo.Changes
			}
			var commitInfoWorkItems []webapi.ResourceRef
			if commitInfo.WorkItems != nil {
				commitInfoWorkItems = *commitInfo.WorkItems
			}

			commitInfoCommitIdShort := commitInfoCommitId[:7]
			commitInfoRemoteRepoUrl := repoUrlPattern.ReplaceAllString(commitInfoRemoteUrl, "$1")

			ci = append(ci, model.CommitInfoSimple{
				Author:           commitInfoAuthor,
				Comment:          commitInfoComment,
				CommentTruncated: commitInfoCommentTruncated,
				CommitId:         commitInfoCommitId,
				CommitIdShort:    commitInfoCommitIdShort,
				Committer:        commitInfoCommitter,
				Changes:          commitInfoChanges,
				Push:             commitInfoPush,
				RemoteUrl:        commitInfoRemoteUrl,
				RemoteRepoUrl:    commitInfoRemoteRepoUrl,
				WorkItems:        commitInfoWorkItems,
			})
			logger.json.Debug("ReturnGitCommitItemSimple", "commitItem.Repository", commitItem.Repository, "*commitInfo.Comment", *commitInfo.Comment)
		}
		commitItems = append(commitItems, model.GitCommitItemSimple{
			Repository: commitItem.Repository,
			CommitInfo: ci,
		})
	}
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
