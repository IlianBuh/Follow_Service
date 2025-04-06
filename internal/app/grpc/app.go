package grpcapp

import (
	"context"
	"fmt"
	"github.com/IlianBuh/Follow_Service/internal/lib/logger/sl"
	grpcfllw "github.com/IlianBuh/Follow_Service/internal/transport/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log     *slog.Logger
	gRPCSrv *grpc.Server
	port    int
}

type Service interface {
	Follow(ctx context.Context, src, target int) error
	Unfollow(ctx context.Context, src, target int) error
	ListFollowers(ctx context.Context, uuid int) ([]int, error)
	ListFollowees(ctx context.Context, uuid int) ([]int, error)
}

func New(
	log *slog.Logger,
	port int,
	srvc Service,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(interface{}) error {
			log.Error("catch panic")
			return nil
		}),
	}

	grpcsrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(
				recoveryOpts...,
			),
			logging.UnaryServerInterceptor(
				logInterceptor(log), loggingOpts...,
			),
		),
	)

	grpcfllw.Register(grpcsrv, srvc)

	return &App{log: log, gRPCSrv: grpcsrv, port: port}
}

// logInterceptor is wrapper for logger to enable convenient my logger for grpc interceptor
func logInterceptor(log *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields)
	})
}

// MustRun starts application and throw panic if error occurred
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("failed to start grpc application: " + err.Error())
	}
}

// Run starts grpc application
func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op))
	log.Info("starting grpc application")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("failed to listen socket", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = a.gRPCSrv.Serve(lis); err != nil {
		log.Error("failed to serve", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stopping grpc application
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.Info("stopping grpc application")

	a.gRPCSrv.GracefulStop()
}
