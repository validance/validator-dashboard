package client

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"sync"
)

const STRIDE = 1000

type CosmosClient interface {
	ValidatorDelegations() ([]staking.DelegationResponse, error)
}

type ValidatorQueryClient interface {
	validatorDelegation(offset uint64, limit uint64) (*staking.QueryValidatorDelegationsResponse, error)
	selfDelegationReward()
	commission() (*distribution.QueryValidatorCommissionResponse, error)
}

type GrantQueryClient interface {
	reward() (*distribution.QueryDelegationTotalRewardsResponse, error)
}

type validatorQueryClient struct {
	validatorAddr           string
	stakingQueryClient      staking.QueryClient
	distributionQueryClient distribution.QueryClient
}

type grantQueryClient struct {
	grantWalletAddr    string
	distributionClient distribution.QueryClient
}

type cosmosClient struct {
	grpcConn             *grpc.ClientConn
	url                  string
	validatorQueryClient ValidatorQueryClient
	grantQueryClient     GrantQueryClient
}

func NewCosmosClient(url string, validatorAddr, grantWalletAddr string) (CosmosClient, error) {
	grpcConn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	stakingClient := staking.NewQueryClient(grpcConn)
	distributionClient := distribution.NewQueryClient(grpcConn)

	vqc := &validatorQueryClient{
		validatorAddr,
		stakingClient,
		distributionClient,
	}

	gqc := &grantQueryClient{
		grantWalletAddr,
		distributionClient,
	}

	return cosmosClient{
		grpcConn,
		url,
		vqc,
		gqc,
	}, nil
}

func (v validatorQueryClient) validatorDelegation(offset uint64, limit uint64) (*staking.QueryValidatorDelegationsResponse, error) {
	res, err := v.stakingQueryClient.ValidatorDelegations(
		context.Background(),
		&staking.QueryValidatorDelegationsRequest{
			ValidatorAddr: v.validatorAddr,
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

func (c cosmosClient) appendDelegationResponses(totalValidatorDelegations *[]staking.DelegationResponse, validatorDelegations staking.DelegationResponses) {
	for _, d := range validatorDelegations {
		*totalValidatorDelegations = append(*totalValidatorDelegations, d)
	}
}

func (c cosmosClient) ValidatorDelegations() ([]staking.DelegationResponse, error) {

	var validatorDelegations []staking.DelegationResponse
	var wg sync.WaitGroup

	// initial fetch to get total data
	res, err := c.validatorQueryClient.validatorDelegation(0, STRIDE)

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
	ch := make(chan staking.DelegationResponses, iterate)

	for i := 0; i < iterate; i++ {
		offset := i + 1
		go func() {
			defer wg.Done()
			vd, err := c.validatorQueryClient.validatorDelegation(uint64(offset*STRIDE), STRIDE)

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

func (v validatorQueryClient) selfDelegationReward() {
}

func (v validatorQueryClient) commission() (*distribution.QueryValidatorCommissionResponse, error) {
	res, err := v.distributionQueryClient.ValidatorCommission(
		context.Background(),
		&distribution.QueryValidatorCommissionRequest{
			ValidatorAddress: v.validatorAddr,
		})

	return res, err
}

func (v validatorQueryClient) validatorIncome() {
	// TODO return sum of reward and commission values

}

func (g grantQueryClient) reward() (*distribution.QueryDelegationTotalRewardsResponse, error) {
	res, err := g.distributionClient.DelegationTotalRewards(context.Background(), &distribution.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: g.grantWalletAddr,
	})

	return res, err

}
