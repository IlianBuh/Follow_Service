package main

import (
	"github.com/IlianBuh/Follow_Service/internal/app"
	"github.com/IlianBuh/Follow_Service/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setUpLogger(cfg.Env)

	log.Info("", slog.Any("cfg", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.StorageURL)

	go application.GRPCApp.MustRun()

	stop := make(chan os.Signal)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("received signal", slog.Any("signal", sign))

	application.GRPCApp.Stop()
}

// setUpLogger returns set logger according to current environment
func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
