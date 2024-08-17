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
	logger.json.Info("RenderDashboard", "NewPATConnection", connection)
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

func (adoClients ADOClients) GetRepositories(ctx context.Context, project string) *[]git.GitRepository {
	responseValue, err := adoClients.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{Project: &project})
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
