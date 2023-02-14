package client

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"sync"
)

const STRIDE = 1000

type CosmosClient struct {
	grpcConn      *grpc.ClientConn
	url           string
	validatorAddr string
}

func NewCosmosClient(url string, validatorAddr string) (*CosmosClient, error) {
	grpcConn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	return &CosmosClient{
		grpcConn, url, validatorAddr,
	}, nil
}

func (c CosmosClient) fetchValidatorDelegation(stakingClient types.QueryClient, offset uint64, limit uint64) (*types.QueryValidatorDelegationsResponse, error) {
	res, err := stakingClient.ValidatorDelegations(
		context.Background(),
		&types.QueryValidatorDelegationsRequest{
			ValidatorAddr: c.validatorAddr,
			Pagination: &query.PageRequest{
				Key:        nil,
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
				Reverse:    true,
			},
		})

	return res, err
}

func (c CosmosClient) appendDelegationResponses(totalValidatorDelegations *[]types.DelegationResponse, validatorDelegations types.DelegationResponses) {
	for _, d := range validatorDelegations {
		*totalValidatorDelegations = append(*totalValidatorDelegations, d)
	}
}

func (c CosmosClient) ValidatorDelegations() ([]types.DelegationResponse, error) {
	stakingClient := types.NewQueryClient(c.grpcConn)

	var validatorDelegations []types.DelegationResponse
	var wg sync.WaitGroup

	// initial fetch to get total data
	res, err := c.fetchValidatorDelegation(stakingClient, 0, STRIDE)

	if err != nil {
		return nil, err
	}

	c.appendDelegationResponses(&validatorDelegations, res.GetDelegationResponses())

	total := res.GetPagination().GetTotal()
	fmt.Printf("total: %d\n", total)

	iterate := int(math.Ceil(float64(total / STRIDE)))
	fmt.Printf("iterate: %d\n", iterate)

	wg.Add(iterate)

	// create iterate-sized buffered channel
	ch := make(chan types.DelegationResponses, iterate)

	for i := 0; i < iterate; i++ {
		offset := i + 1
		go func() {
			defer wg.Done()
			vd, err := c.fetchValidatorDelegation(stakingClient, uint64(offset*STRIDE), STRIDE)

			if err != nil {
				fmt.Println(err)
				return
			}

			ch <- vd.GetDelegationResponses()
		}()
	}

	// wait until every goroutine finish pushing data into channel
	wg.Wait()

	// close buffered channel
	close(ch)

	// pop data from buffered channel
	for vd := range ch {
		c.appendDelegationResponses(&validatorDelegations, vd)
	}

	return validatorDelegations, err
}

func (c CosmosClient) fetchReward() {

}

func (c CosmosClient) fetchCommission() {

}

func (c CosmosClient) ValidatorIncome() {
	// TODO return sum of reward and commission values
	c.fetchReward()
	c.fetchCommission()
}
