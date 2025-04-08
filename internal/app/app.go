package app

import (
	"fmt"
	grpcapp "github.com/IlianBuh/Follow_Service/internal/app/grpc"
	grpclient "github.com/IlianBuh/Follow_Service/internal/clients/grpc"
	"github.com/IlianBuh/Follow_Service/internal/service/follow"
	"github.com/IlianBuh/Follow_Service/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	usrInfoPort int,
	storagePath string,
	retryCount int,
	timeout time.Duration,
) *App {
	st, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	cl, err := grpclient.New(log, fmt.Sprintf(":%d", usrInfoPort), retryCount, timeout)
	if err != nil {
		panic(err)
	}
	fl := follow.New(log, st, st, st, cl)

	application := grpcapp.New(log, port, fl)

	return &App{
		GRPCApp: application,
	}
}
