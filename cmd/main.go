package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/application"
	"github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/config"
	mlog "github.com/yuanyu90221/uniplile-authentication-with-linkedin/internal/logger"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		},
	))

	rootContext := context.WithValue(context.Background(), mlog.CtxKey{}, logger)
	config.Init(rootContext)
	app := application.New(rootContext, config.AppConfig)
	ctx, cancel := signal.NotifyContext(rootContext, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		logger.Error("failed to start app", "error", err)
	}
}
