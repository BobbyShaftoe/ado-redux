package main

import (
	"HTTP_Sever/handlers"
	"HTTP_Sever/helpers/ado"
	"HTTP_Sever/helpers/config"
	"HTTP_Sever/model"
	"context"
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
	envVars := config.New()

	adoCtx := context.Background()

	userValidated := ado.ValidateUser(envVars.USER, handlers.NewADOClients(adoCtx).ListUsers(adoCtx, envVars.PROJECT))

	globalState := &model.GlobalState{
		User:           envVars.USER,
		UserValidated:  userValidated,
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

	logger.json.Info("main", "status", "Starting server at http://localhost:8080")
	fatal(http.ListenAndServe(":8080", nil))
}
