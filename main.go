package main

import (
	"HTTP_Sever/handlers"
	"HTTP_Sever/helpers/ado"
	"context"
	"fmt"
	"github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	DBPass string
	PAT    string
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

	adoClientInfo := handlers.GetADOClientInfo("https://dev.azure.com/Kingsizenix", envVars.PAT)
	fmt.Println(adoClientInfo)

	adoConnection := handlers.NewPATConnection(adoClientInfo)
	adoCtx := context.Background()
	//coreClient := handlers.NewADOClient(adoCtx, adoConnection)
	gitClient := handlers.NewGitClient(adoCtx, adoConnection)

	//responseValue, err := coreClient.GetProjects(adoCtx, core.GetProjectsArgs{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(ado.ReturnProjects(responseValue))

	responseValue2, err := gitClient.GetRepositories(adoCtx, git.GetRepositoriesArgs{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ado.ReturnGitRepos(responseValue2))

	tMap := handlers.NewTemplMap()

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", templ.Handler(handlers.RenderRouteTempl(*tMap, "index")))
	http.Handle("/hello", templ.Handler(handlers.RenderRouteTempl(*tMap, "hello")))
	http.Handle("/dashboard", templ.Handler(handlers.RenderRouteTempl(*tMap, "dashboard")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (envvars *EnvVars) getEnv() {
	envvars.DBPass = os.Getenv("DB_PASS")
	envvars.PAT = os.Getenv("AZURE_TOKEN")
	if envvars.DBPass == "" {
		log.Fatal("DB_PASS environment variable not set")
	}
	if envvars.PAT == "" {
		log.Fatal("AZURE_TOKEN environment variable not set")
	}
}
