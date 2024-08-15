package main

import (
	"HTTP_Sever/handlers"
	"HTTP_Sever/model"
	"fmt"
	"github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"net/http"
	"os"
	"sync"
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

type EnvVars struct {
	DBPass       string
	PAT          string
	ORGANIZATION string
	PROJECT      string
}

const (
	host   = "localhost"
	port   = 5432
	user   = "go_user"
	dbname = "go_database"
)

var (
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

	logger.json.Info("main", "globalState", globalState)

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", templ.Handler(handlers.RenderIndex(globalState)))

	http.Handle("/hello", templ.Handler(handlers.RenderHello(globalState)))

	http.Handle("/dashboard", templ.Handler(handlers.RenderDashboard(globalState)))

	http.HandleFunc("/dashboard-update", handlers.RenderDashboardUpdateProject(globalState))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	logger.json.Info("main", "msg", "Starting server at http://localhost:8080")
	fatal(http.ListenAndServe(":8080", nil))
}

func (envVars *EnvVars) getEnv() {
	envVars.DBPass = os.Getenv("DB_PASS")
	envVars.PAT = os.Getenv("AZURE_TOKEN")
	envVars.ORGANIZATION = os.Getenv("ADO_ORG")
	envVars.PROJECT = os.Getenv("ADO_DEFAULT_PROJECT")

	if envVars.DBPass == "" {
		fatal("DB_PASS environment variable not set")
	}
	if envVars.PAT == "" {
		fatal("AZURE_TOKEN environment variable not set")
	}
	if envVars.ORGANIZATION == "" {
		fatal("ADO_ORG environment variable not set")
	}
	if envVars.PROJECT == "" {
		fatal("ADO_DEFAULT_PROJECT environment variable not set")
	}
}
