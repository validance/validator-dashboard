package main

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConn, err := grpc.Dial(
		"localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(err)
	}

	stakingClient := types.NewQueryClient(grpcConn)

	res, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&types.QueryValidatorDelegationsRequest{
			ValidatorAddr: "cosmosvaloper1v78emy9d2xe3tj974l7tmn2whca2nh9zp7s0u9",
			Pagination: &query.PageRequest{
				Key:        nil,
				Offset:     0,
				Limit:      100,
				CountTotal: true,
				Reverse:    false,
			},
		})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
