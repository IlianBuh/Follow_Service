package tests

import (
	followv1 "github.com/IlianBuh/Follow_Protobuf/gen/go"
	"github.com/IlianBuh/Follow_Service/tests/suite"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestFollowUnfollowHappy(t *testing.T) {
	ctx, st := suite.New(t)

	const seed = int64(1)
	rand := rand.New(rand.NewSource(seed))

	src, target := rand.Uint32()&uint32(1<<31-1), rand.Uint32()&uint32(1<<31-1)

	_, err := st.Client.Follow(
		ctx,
		&followv1.FollowRequest{
			Src:    int32(src),
			Target: int32(target),
		},
	)
	require.NoError(t, err)

	_, err = st.Client.Unfollow(
		ctx,
		&followv1.UnfollowRequest{
			Src:    int32(src),
			Target: int32(target),
		},
	)
	require.NoError(t, err)
}
