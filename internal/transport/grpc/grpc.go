package grpcfllw

import (
	"context"
	"fmt"
	followv1 "github.com/IlianBuh/Follow_Protobuf/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	Follow(ctx context.Context, src, target int) error
	Unfollow(ctx context.Context, src, target int) error
	ListFollowers(ctx context.Context, uuid int) ([]int, error)
	ListFollowees(ctx context.Context, uuid int) ([]int, error)
}
type serverAPI struct {
	fllw Service
	followv1.UnimplementedFollowServer
}

// Register registers handlers on grpc server
func Register(grpcsrv *grpc.Server, fllw Service) {
	followv1.RegisterFollowServer(grpcsrv, &serverAPI{fllw: fllw})
}

// Follow is API-handler for Follow method
func (s *serverAPI) Follow(
	ctx context.Context,
	req *followv1.FollowRequest,
) (*followv1.FollowResponse, error) {
	pars := int32ToInt(req.GetSrc(), req.GetTarget())

	if err := validateIntValues(pars[0], pars[1]); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.fllw.Follow(ctx, pars[0], pars[1])
	if err != nil {
		// TODO : handle error when users are not found

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &followv1.FollowResponse{}, nil
}

// Unfollow is API-handler for Unfollow method
func (s *serverAPI) Unfollow(
	ctx context.Context,
	req *followv1.UnfollowRequest,
) (*followv1.UnfollowResponse, error) {
	pars := int32ToInt(req.GetSrc(), req.GetTarget())

	if err := validateIntValues(pars[0], pars[1]); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.fllw.Unfollow(ctx, pars[0], pars[1])
	if err != nil {
		// TODO : handle error when users are not found

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &followv1.UnfollowResponse{}, nil
}

// ListFollowers is API-handler for ListFollowers method
func (s *serverAPI) ListFollowers(
	ctx context.Context,
	req *followv1.ListFollowersRequest,
) (*followv1.ListFollowersResponse, error) {
	pars := int32ToInt(req.GetUuid())

	if err := validateIntValues(pars[0]); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	uuids, err := s.fllw.ListFollowers(ctx, pars[0])
	if err != nil {
		// TODO : handle error when users are not found

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &followv1.ListFollowersResponse{Uuids: intToInt32(uuids...)}, nil
}

// ListFollowees is API-handler for ListFollowees method
func (s *serverAPI) ListFollowees(
	ctx context.Context,
	req *followv1.ListFolloweesRequest,
) (*followv1.ListFolloweesResponse, error) {
	pars := int32ToInt(req.GetUuid())

	if err := validateIntValues(pars[0]); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	uuids, err := s.fllw.ListFollowees(ctx, pars[0])
	if err != nil {
		// TODO : handle error when users are not found

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &followv1.ListFolloweesResponse{Uuids: intToInt32(uuids...)}, nil
}

// validateIntValues validates value to be non-negative
func validateIntValues(vals ...int) error {
	for _, v := range vals {
		if v < 0 {
			return fmt.Errorf("uuid can't be negative")
		}
	}

	return nil
}

// int32ToInt converts list of int32 values to slice of int
func int32ToInt(vals ...int32) []int {
	res := make([]int, len(vals))

	for i := range vals {
		res[i] = int(vals[i])
	}

	return res
}

// intToInt32 converts list of int values to slice of int32
func intToInt32(vals ...int) []int32 {
	res := make([]int32, len(vals))

	for i := range vals {
		res[i] = int32(vals[i])
	}

	return res
}
