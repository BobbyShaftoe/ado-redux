package handlers

import (
	"HTTP_Sever/helpers/config"
	"HTTP_Sever/model"
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"log"
	"sort"
	"sync"
)

type ADOClients struct {
	coreClient  core.Client
	gitClient   git.Client
	graphClient graph.Client
	buildClient build.Client
}

type ADORequests interface {
	GetProjects(ctx context.Context, coreClient core.Client) *core.GetProjectsResponseValue
	GetRepositories(ctx context.Context, gitClient git.Client) *[]git.GitRepository
}

func GetADOClientInfo() model.ADOConnectionInfo {
	adoConnectionInfo := model.ADOConnectionInfo{
		ConnectionUrl: "https://dev.azure.com/" + config.New().Organization,
		ConnectionPAT: config.New().Pat,
	}
	return adoConnectionInfo
}

func NewPATConnection() *azuredevops.Connection {
	adoClientInfo := GetADOClientInfo()
	connection := azuredevops.NewPatConnection(adoClientInfo.ConnectionUrl, adoClientInfo.ConnectionPAT)
	logger.json.Debug("RenderDashboard", "NewPATConnection", connection)
	return connection
}

func NewADOClients(ctx context.Context) *ADOClients {
	patConnection := NewPATConnection()

	coreClient, err := core.NewClient(ctx, patConnection)
	if err != nil {
		fatal(err)
	}
	gitClient, err := git.NewClient(ctx, patConnection)
	if err != nil {
		fatal(err)
	}
	buildClient, err := build.NewClient(ctx, patConnection)
	if err != nil {
		fatal(err)
	}
	graphClient, err := graph.NewClient(ctx, patConnection)
	return &ADOClients{
		coreClient:  coreClient,
		gitClient:   gitClient,
		graphClient: graphClient,
		buildClient: buildClient,
	}

}

func (adoClients ADOClients) GetProjects(ctx context.Context) *core.GetProjectsResponseValue {
	responseValue, err := adoClients.coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	if err != nil {
		log.Fatal(err)
	}
	return responseValue
}

func (adoClients ADOClients) GetRepositories(ctx context.Context, globalState *model.GlobalState) *[]git.GitRepository {
	logger.json.Debug("GetRepositories", "project", globalState.CurrentProject)
	responseValue, err := adoClients.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{Project: &globalState.CurrentProject})
	if err != nil {
		log.Fatal(err)
	}
	return responseValue
}

func ReturnGitCommitCriteria(globalState *model.GlobalState) *model.GitCommitsCriteria {
	return &model.GitCommitsCriteria{
		//RepositoryId: "",
		Author:      globalState.User,
		User:        globalState.User,
		FromDate:    "6/14/2018 12:00:00 AM",
		Version:     "",
		VersionType: "",
		Skip:        0,
		Top:         10,
	}

}

func (adoClients ADOClients) GetCommits(ctx context.Context, gcc *model.GitCommitsCriteria, globalState *model.GlobalState, repositories []string) []model.GitCommitItem {
	var True = true
	var False = false
	var allCommits []model.GitCommitItem

	//var gitVersionDescriptor = git.GitVersionDescriptor{
	//	Version:        &gitCommitsCriteria.Version,
	//	VersionOptions: nil,
	//	VersionType:    (*git.GitVersionType)(&gitCommitsCriteria.VersionType),
	//}

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, repo := range repositories {
		wg.Add(1)

		go func(repo string) {
			defer wg.Done()

			gitCommitsCriteria := model.GitCommitsCriteria{
				Author:      gcc.Author,
				User:        gcc.User,
				FromDate:    gcc.FromDate,
				Version:     gcc.Version,
				VersionType: gcc.VersionType,
				Skip:        gcc.Skip,
				Top:         gcc.Top,
			}

			gitCommitsCriteria.RepositoryId = repo

			responseValue, err := adoClients.gitClient.GetCommits(ctx, git.GetCommitsArgs{
				Project:      &globalState.CurrentProject,
				RepositoryId: &gitCommitsCriteria.RepositoryId,
				SearchCriteria: &git.GitQueryCommitsCriteria{
					Skip: &gitCommitsCriteria.Skip,
					Top:  &gitCommitsCriteria.Top,
					//Author:              &gitCommitsCriteria.Author,
					CompareVersion:      nil,
					ExcludeDeletes:      &True,
					FromCommitId:        nil,
					FromDate:            &gitCommitsCriteria.FromDate,
					HistoryMode:         nil,
					Ids:                 nil,
					IncludeLinks:        &True,
					IncludePushData:     &True,
					IncludeUserImageUrl: nil,
					IncludeWorkItems:    &True,
					ItemPath:            nil,
					ItemVersion:         nil,
					//ItemVersion:            &gitVersionDescriptor,
					ShowOldestCommitsFirst: &False,
					ToCommitId:             nil,
					ToDate:                 nil,
					User:                   &gitCommitsCriteria.User, // User is the commit author
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			responseCommits := responseValue
			if len(*responseCommits) == 0 {
				return
			}

			logger.json.Debug("GetCommits", "gitCommitsCriteria", gitCommitsCriteria)

			mu.Lock()
			allCommits = append(allCommits, model.GitCommitItem{
				Repository: repo,
				CommitInfo: *responseCommits,
			})
			mu.Unlock()

		}(repo)

	}
	wg.Wait()

	sort.Slice(allCommits, func(i, j int) bool {
		return allCommits[i].Repository < allCommits[j].Repository
	})
	return allCommits
}

func (adoClients ADOClients) GetPush(ctx context.Context, repository string, pushId int, globalState *model.GlobalState) *git.GitPush {
	responseValue, err := adoClients.gitClient.GetPush(ctx, git.GetPushArgs{
		Project:      &globalState.CurrentProject,
		RepositoryId: &repository,
		PushId:       &pushId,
	})
	if err != nil {
		log.Fatal(err)
	}
	return responseValue
}

// The ListBuilds function in the ADOClients struct is designed to retrieve a list of builds from Azure DevOps for a given set of repositories.
// It takes three parameters: a context.Context object, a globalState object of type *model.GlobalState, and a slice of repository names.
// https://medium.com/lyonas/go-type-casting-starter-guide-a9c1811670c5
func (adoClients ADOClients) ListBuilds(ctx context.Context, globalState *model.GlobalState, repositories []model.GitRepo) []build.Build {
	buildsList := make([]build.Build, 0)

	top := 1
	var repositoryType = "TfsGit"
	var queryOrder build.BuildQueryOrder = "descending"

	for _, repo := range repositories {
		repoId := repo.Id.String()
		logger.json.Debug("ListBuilds", "repoId", repoId)

		var responseValue, err = adoClients.buildClient.GetBuilds(ctx, build.GetBuildsArgs{
			Project:                &globalState.CurrentProject,
			Definitions:            nil,
			Queues:                 nil,
			BuildNumber:            nil,
			MinTime:                nil,
			MaxTime:                nil,
			RequestedFor:           nil,
			ReasonFilter:           nil,
			StatusFilter:           nil,
			ResultFilter:           nil,
			TagFilters:             nil,
			Properties:             nil,
			Top:                    &top,
			ContinuationToken:      nil,
			MaxBuildsPerDefinition: nil,
			DeletedFilter:          nil,
			QueryOrder:             &queryOrder,
			BranchName:             nil,
			BuildIds:               nil,
			RepositoryId:           &repoId,
			RepositoryType:         &repositoryType,
		})
		if err != nil {
			continue
		}
		builds := responseValue
		buildsList = append(buildsList, builds.Value...)
	}
	logger.json.Debug("ListBuilds", "buildsList", buildsList)
	return buildsList
}

// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/latest/get?view=azure-devops-rest-7.2#buildrepository
func (adoClients ADOClients) GetLatestBuildsFromBuilds(ctx context.Context, globalState *model.GlobalState, buildList []build.Build) []build.Build {
	latestBuilds := make([]build.Build, 0)
	for _, buildItem := range buildList {
		//definitionId := strconv.Itoa(*buildItem.Definition.Id)
		responseValue, err := adoClients.buildClient.GetLatestBuild(ctx, build.GetLatestBuildArgs{
			Project:    &globalState.CurrentProject,
			Definition: buildItem.Definition.Name,
			//Definition: &definitionId,
			BranchName: buildItem.SourceBranch,
		})

		if err != nil {
			logger.json.Info("GetLatestBuildsFromBuilds", "err", err, "repo", buildItem.Definition.Name)
			continue
		}

		currentBuildItem := *responseValue
		latestBuilds = append(latestBuilds, currentBuildItem)

	}
	return latestBuilds
}

func (adoClients ADOClients) GetLatestBuildsFromRepositories(ctx context.Context, globalState *model.GlobalState, repositories []string) []build.Build {
	latestBuilds := make([]build.Build, 0)
	for _, repo := range repositories {
		logger.json.Info("GetLatestBuildsFromRepositories", "repo", repo)
		responseValue, err := adoClients.buildClient.GetLatestBuild(ctx, build.GetLatestBuildArgs{
			Project:    &globalState.CurrentProject,
			Definition: &repo,
			BranchName: nil,
		})
		if err != nil {
			//logger.json.Error("GetLatestBuildsFromRepositories", "err", err)
			continue
		}
		buildItem := *responseValue
		latestBuilds = append(latestBuilds, buildItem)
		//logger.json.Info("GetLatestBuildsFromRepositories", "latestBuild", responseValue)
	}
	return latestBuilds
}

func (adoClients ADOClients) ListUsers(ctx context.Context, project string) *[]graph.GraphUser {
	subjectTypes := &[]string{"msa", "aad"}
	//scopeDescriptor := project
	var graphUsers []graph.GraphUser

	responseValue, err := adoClients.graphClient.ListUsers(
		ctx, graph.ListUsersArgs{
			ContinuationToken: nil,
			//ScopeDescriptor:   &scopeDescriptor,
			SubjectTypes: subjectTypes,
		},
	)
	if err != nil {
		fatal(err)
	}

	graphUsers = append(graphUsers, *responseValue.GraphUsers...)
	ct := *responseValue.ContinuationToken

	if ct[0] == "" {
		return &graphUsers
	}

	for ct[0] != "" {
		if err != nil {
			fatal(err)
		}

		listUserArgs := graph.ListUsersArgs{
			ContinuationToken: &ct[0],
		}

		responseValue, err = adoClients.graphClient.ListUsers(
			ctx, listUserArgs,
		)
		if err != nil {
			fatal(err)
		}

		graphUsers = append(graphUsers, *responseValue.GraphUsers...)
		ct = *responseValue.ContinuationToken
	}
	return &graphUsers
}
