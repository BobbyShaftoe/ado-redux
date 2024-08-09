package handlers

import (
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"log"
)

type ADOClientInfo struct {
	organizationUrl     string
	personalAccessToken string
}

type SimpleADOConnection struct{}

func GetADOClientInfo(organizationUrl string, personalAccessToken string) ADOClientInfo {
	adoClientInfo := ADOClientInfo{
		organizationUrl:     organizationUrl,
		personalAccessToken: personalAccessToken,
	}
	return adoClientInfo
}

func (c SimpleADOConnection) NewPATConnection(organizationUrl string, personalAccessToken string) *azuredevops.Connection {
	adoClientInfo := GetADOClientInfo(organizationUrl, personalAccessToken)
	connection := azuredevops.NewPatConnection(adoClientInfo.organizationUrl, adoClientInfo.personalAccessToken)
	return connection
}

func NewPATConnection(adoClientInfo ADOClientInfo) *azuredevops.Connection {
	connection := azuredevops.NewPatConnection(adoClientInfo.organizationUrl, adoClientInfo.personalAccessToken)
	return connection
}

func NewADOClient(ctx context.Context, connection *azuredevops.Connection) core.Client {
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}
	return coreClient
}

func NewGitClient(ctx context.Context, connection *azuredevops.Connection) git.Client {
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}
	return gitClient
}
