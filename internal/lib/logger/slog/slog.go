package slog

import (
	"applicationDesignTest/internal/consts"
	"log/slog"
	"os"
	"strings"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch {
	case strings.EqualFold(env, consts.EnvLocal):
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case strings.EqualFold(env, consts.EnvDevelop):
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case strings.EqualFold(env, consts.EnvProduction):
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
