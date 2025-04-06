package tests

import (
	followv1 "github.com/IlianBuh/Follow_Protobuf/gen/go"
	"github.com/IlianBuh/Follow_Service/tests/suite"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestListFollowers(t *testing.T) {
	ctx, st := suite.New(t)

	rand := rand.New(rand.NewSource(time.Now().Unix()))

	uuid := randUUID(rand)
	followers := randomInt32Slice(10, rand)
	for _, v := range followers {
		_, err := st.Client.Follow(
			ctx,
			&followv1.FollowRequest{
				Src:    v,
				Target: uuid,
			},
		)
		require.NoError(t, err)
	}

	res, err := st.Client.ListFollowers(
		ctx,
		&followv1.ListFollowersRequest{
			Uuid: uuid,
		},
	)
	require.NoError(t, err)
	require.Equal(t, followers, res.GetUuids())
}

func randomInt32Slice(size int, rand *rand.Rand) []int32 {
	res := make([]int32, size)

	size--
	for size >= 0 {
		res[size] = randUUID(rand)
		size--
	}

	return res
}
func randUUID(rand *rand.Rand) int32 {
	return int32(rand.Uint32() & (1<<31 - 1))
}
