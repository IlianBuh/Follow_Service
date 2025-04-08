package grpclient

import (
	"context"
	"fmt"
	"github.com/IlianBuh/Follow_Service/internal/lib/mappers"
	userinfov1 "github.com/IlianBuh/SSO_Protobuf/gen/go/userinfo"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func (c *Client) CheckUsers(ctx context.Context, uuids []int) (bool, error) {
	const op = "grpclient.CheckUsers"

	res, err := c.gRPClient.UsersExist(
		ctx,
		&userinfov1.UsersExistRequest{
			Uuid: mappers.IntToInt32(uuids...),
		},
	)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res.GetExist(), nil
}
