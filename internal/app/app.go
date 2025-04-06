package app

import (
	grpcapp "github.com/IlianBuh/Follow_Service/internal/app/grpc"
	"github.com/IlianBuh/Follow_Service/internal/service/follow"
	"github.com/IlianBuh/Follow_Service/internal/storage/sqlite"
	"log/slog"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	storagePath string,
) *App {
	st, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	fl := follow.New(log, st, st, st)

	application := grpcapp.New(log, port, fl)

	return &App{
		GRPCApp: application,
	}
}
