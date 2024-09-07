package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type EnvVars struct {
	Pat             string
	Organization    string
	Project         string
	User            string
	AdditionalUsers []string
}

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

func New() *EnvVars {
	envVars := &EnvVars{}
	envVars.getEnv()
	return envVars
}

func (EnvVars *EnvVars) getEnv() {
	EnvVars.Pat = os.Getenv("AZURE_TOKEN")
	EnvVars.Organization = os.Getenv("ADO_ORG")
	EnvVars.Project = os.Getenv("ADO_DEFAULT_PROJECT")
	EnvVars.User = os.Getenv("ADO_DEFAULT_USER")

	if additionalUsers := os.Getenv("ADO_ADDITIONAL_USERS"); additionalUsers != "" {
		EnvVars.AdditionalUsers = strings.Split(additionalUsers, ",")
	} else {
		EnvVars.AdditionalUsers = make([]string, 0)
	}

	if EnvVars.Pat == "" {
		fatal("AZURE_TOKEN environment variable not set")
	}
	if EnvVars.Organization == "" {
		fatal("ADO_ORG environment variable not set")
	}
	if EnvVars.Project == "" {
		fatal("ADO_DEFAULT_PROJECT environment variable not set")
	}
	if EnvVars.User == "" {
		fatal("ADO_DEFAULT_USER environment variable not set")
	}
}
