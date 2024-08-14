package handlers

import (
	"HTTP_Sever/model"
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"log"
	"os"
)

type ADOClients struct {
	coreClient core.Client
	gitClient  git.Client
}

type ADORequests interface {
	GetProjects(ctx context.Context, coreClient core.Client) *core.GetProjectsResponseValue
	GetRepositories(ctx context.Context, gitClient git.Client) *[]git.GitRepository
}

func GetADOClientInfo() model.ADOConnectionInfo {
	adoConnectionInfo := model.ADOConnectionInfo{
		ConnectionUrl: "https://dev.azure.com/" + os.Getenv("ADO_ORG"),
		ConnectionPAT: os.Getenv("AZURE_TOKEN"),
	}
	return adoConnectionInfo
}

func NewPATConnection() *azuredevops.Connection {
	adoClientInfo := GetADOClientInfo()
	connection := azuredevops.NewPatConnection(adoClientInfo.ConnectionUrl, adoClientInfo.ConnectionPAT)
	return connection
}

func NewADOClients(ctx context.Context) *ADOClients {
	patConnection := NewPATConnection()

	coreClient, err := core.NewClient(ctx, patConnection)
	if err != nil {
		log.Fatal(err)
	}
	gitClient, err := git.NewClient(ctx, patConnection)
	if err != nil {
		log.Fatal(err)
	}
	return &ADOClients{
		coreClient: coreClient,
		gitClient:  gitClient,
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
