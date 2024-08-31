package handlers

import (
	"HTTP_Sever/helpers/config"
	"HTTP_Sever/model"
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"log"
	"sync"
)

type ADOClients struct {
	coreClient  core.Client
	gitClient   git.Client
	graphClient graph.Client
}

type ADORequests interface {
	GetProjects(ctx context.Context, coreClient core.Client) *core.GetProjectsResponseValue
	GetRepositories(ctx context.Context, gitClient git.Client) *[]git.GitRepository
}

func GetADOClientInfo() model.ADOConnectionInfo {
	adoConnectionInfo := model.ADOConnectionInfo{
		ConnectionUrl: "https://dev.azure.com/" + config.New().ORGANIZATION,
		ConnectionPAT: config.New().PAT,
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
	if err != nil {
		fatal(err)
	}
	graphClient, err := graph.NewClient(ctx, patConnection)
	return &ADOClients{
		coreClient:  coreClient,
		gitClient:   gitClient,
		graphClient: graphClient,
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
		Top:         3,
	}

}

func (adoClients ADOClients) GetCommits(ctx context.Context, gcc *model.GitCommitsCriteria, globalState *model.GlobalState, repositories []string) []model.GitCommitItem {
	var True = true
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
					Skip:                &gitCommitsCriteria.Skip,
					Top:                 &gitCommitsCriteria.Top,
					Author:              &gitCommitsCriteria.Author,
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
					ShowOldestCommitsFirst: nil,
					ToCommitId:             nil,
					ToDate:                 nil,
					User:                   nil,
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
