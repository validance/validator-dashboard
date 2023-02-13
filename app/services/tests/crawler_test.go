package tests

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"testing"
)

func Test_client(t *testing.T) {
	grpcConn, err := grpc.Dial(
		"https://grpc.cosmos.silknodes.io",
	)

	if err != nil {
		panic(err)
	}

	stakingClient := types.NewQueryClient(grpcConn)

	res, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&types.QueryValidatorDelegationsRequest{
			ValidatorAddr: "",
			Pagination:    nil,
		})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)

}
