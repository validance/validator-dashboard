package cosmos

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcConn      *grpc.ClientConn
	url           string
	validatorAddr string
}

func NewClient(url string, validatorAddr string) (*Client, error) {
	grpcConn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		grpcConn, url, validatorAddr,
	}, nil
}

func (c Client) ValidatorDelegations() (*types.QueryValidatorDelegationsResponse, error) {
	stakingClient := types.NewQueryClient(c.grpcConn)

	res, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&types.QueryValidatorDelegationsRequest{
			ValidatorAddr: c.validatorAddr,
			Pagination: &query.PageRequest{
				Key:        nil,
				Offset:     0,
				Limit:      100,
				CountTotal: true,
				Reverse:    false,
			},
		})

	return res, err
}
