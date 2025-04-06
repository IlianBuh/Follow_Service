package suite

import (
	"context"
	followv1 "github.com/IlianBuh/Follow_Protobuf/gen/go"
	"github.com/IlianBuh/Follow_Service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Client followv1.FollowClient
	Cfg    *config.Config
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("C:\\Social-media\\Follow-service\\Service\\config\\config.yml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	addr := net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port))
	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}

	client := followv1.NewFollowClient(cc)
	return ctx, &Suite{
		T:      t,
		Cfg:    cfg,
		Client: client,
	}
}
