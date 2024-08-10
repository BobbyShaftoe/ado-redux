package main

import (
	"HTTP_Sever/handlers"
	"HTTP_Sever/helpers/ado"
	"HTTP_Sever/model"
	"context"
	"fmt"
	"github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"log"
	"net/http"
	"os"
	"sync"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type EnvVars struct {
	DBPass  string
	PAT     string
	PROJECT string
}

type User struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	LocationID   int
}

type UserQuery struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	Location     string
}

const (
	host   = "localhost"
	port   = 5432
	user   = "go_user"
	dbname = "go_database"
)

var (
	posts   = make(map[int]Post)
	nextID  = 1
	postsMu sync.Mutex
)

func main() {
	envVars := EnvVars{}
	envVars.getEnv()

	adoClientInfo := handlers.GetADOClientInfo("https://dev.azure.com/"+envVars.PROJECT, envVars.PAT)
	fmt.Println(adoClientInfo)

	adoConnection := handlers.NewPATConnection(adoClientInfo)
	adoCtx := context.Background()
	coreClient := handlers.NewADOClient(adoCtx, adoConnection)
	gitClient := handlers.NewGitClient(adoCtx, adoConnection)

	responseValue, err := coreClient.GetProjects(adoCtx, core.GetProjectsArgs{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ado.ReturnProjects(responseValue))

	responseValue2, err := gitClient.GetRepositories(adoCtx, git.GetRepositoriesArgs{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ado.ReturnGitRepos(responseValue2))

	dashboardData := model.DashboardData{
		Projects: ado.ReturnProjects(responseValue),
		Repos:    ado.ReturnGitRepos(responseValue2),
	}

	//helloData := handlers.HelloData{
	//	Name: "Nick",
	//}

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", templ.Handler(handlers.RenderIndex()))

	http.Handle("/hello", templ.Handler(handlers.RenderHello("Nick")))

	http.Handle("/dashboard", templ.Handler(handlers.RenderDashboard(dashboardData)))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (envvars *EnvVars) getEnv() {
	envvars.DBPass = os.Getenv("DB_PASS")
	envvars.PAT = os.Getenv("AZURE_TOKEN")
	envvars.PROJECT = os.Getenv("ADO_PROJECT")

	if envvars.DBPass == "" {
		log.Fatal("DB_PASS environment variable not set")
	}
	if envvars.PAT == "" {
		log.Fatal("AZURE_TOKEN environment variable not set")
	}
	if envvars.PROJECT == "" {
		log.Fatal("PROJECT environment variable not set")
	}
}
