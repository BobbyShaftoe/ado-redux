package main

import (
	"HTTP_Sever/handlers"
	"HTTP_Sever/model"
	"fmt"
	"github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	DBPass       string
	PAT          string
	ORGANIZATION string
	PROJECT      string
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

	globalState := &model.GlobalState{
		User:           "Nick",
		Projects:       nil,
		CurrentProject: envVars.PROJECT,
	}

	fmt.Println("GLOBAL STATE")
	fmt.Println(globalState)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", templ.Handler(handlers.RenderIndex(globalState)))

	http.Handle("/hello", templ.Handler(handlers.RenderHello(globalState)))

	http.Handle("/dashboard", templ.Handler(handlers.RenderDashboard(globalState)))

	http.HandleFunc("/dashboard-update", handlers.RenderDashboardUpdateProject(globalState))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (envVars *EnvVars) getEnv() {
	envVars.DBPass = os.Getenv("DB_PASS")
	envVars.PAT = os.Getenv("AZURE_TOKEN")
	envVars.ORGANIZATION = os.Getenv("ADO_ORG")
	envVars.PROJECT = os.Getenv("ADO_DEFAULT_PROJECT")

	if envVars.DBPass == "" {
		log.Fatal("DB_PASS environment variable not set")
	}
	if envVars.PAT == "" {
		log.Fatal("AZURE_TOKEN environment variable not set")
	}
	if envVars.ORGANIZATION == "" {
		log.Fatal("ADO_ORG environment variable not set")
	}
	if envVars.PROJECT == "" {
		log.Fatal("ADO_DEFAULT_PROJECT environment variable not set")
	}
}
