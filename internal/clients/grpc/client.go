package grpclient

import (
	userinfov1 "github.com/IlianBuh/SSO_Protobuf/gen/go/userinfo"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log/slog"
	"time"
)

type Client struct {
	log       *slog.Logger
	gRPClient userinfov1.UserInfoClient
}

func New(
	log *slog.Logger,
	addr string,
	retryCounts int,
	timeout time.Duration,
) (*Client, error) {

	retryOpts := []retry.CallOption{
		retry.WithMax(uint(retryCounts)),
		retry.WithCodes(codes.Aborted, codes.NotFound, codes.DeadlineExceeded),
		retry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.NewClient(
		addr,
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, err
	}

	gRPClient := userinfov1.NewUserInfoClient(cc)

	return &Client{
		log:       log,
		gRPClient: gRPClient,
	}, nil
}
