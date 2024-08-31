package handlers

import (
	"fmt"
	"log/slog"
	"os"
)

type localLogger struct {
	json  *slog.Logger
	jsonf *slog.Logger
}

var handlerOpts = &slog.HandlerOptions{
	Level: slog.LevelInfo,
}

var logger = &localLogger{
	json:  slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts)),
	jsonf: slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts)),
}

func fatal(v ...any) {
	logger.json.Error("main", "err", fmt.Sprint(v...))
	os.Exit(1)
}
