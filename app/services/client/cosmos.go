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
	"math/big"
	"sync"
)

const stride = 1000

type grantCommission struct {
	delegatorAddr string
	commission    *distribution.QueryDelegationRewardsResponse
}

type cosmosQuerier interface {
	ValidatorDelegations() (map[string]staking.DelegationResponse, error)
	ValidatorIncome() (*big.Int, error)
	GrantRewards() (map[string]*distribution.QueryDelegationRewardsResponse, error)
}

type validatorQuerier interface {
	validatorDelegation(offset uint64, limit uint64) (*staking.QueryValidatorDelegationsResponse, error)
	selfDelegationReward() (*distribution.QueryDelegationRewardsResponse, error)
	commission() (*distribution.QueryValidatorCommissionResponse, error)
}

type grantQuerier interface {
	rewards() (map[string]*distribution.QueryDelegationRewardsResponse, error)
}

type validatorQueryClient struct {
	validatorOperatorAddr   string
	validatorAddr           string
	stakingQueryClient      staking.QueryClient
	distributionQueryClient distribution.QueryClient
}

type grantQueryClient struct {
	validatorOperatorAddr string
	grantWalletAddr       []string
	distributionClient    distribution.QueryClient
}

type cosmosClient struct {
	grpcConn             *grpc.ClientConn
	url                  string
	validatorQueryClient validatorQuerier
	grantQueryClient     grantQuerier
}

// NewCosmosClient create query client for cosmos app-chains
func NewCosmosClient(url string, validatorOperatorAddr string, validatorAddr string, grantWalletAddr ...string) (cosmosQuerier, error) {
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
		validatorOperatorAddr,
		validatorAddr,
		stakingClient,
		distributionClient,
	}

	gqc := &grantQueryClient{
		validatorOperatorAddr,
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
			ValidatorAddr: v.validatorOperatorAddr,
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

func (c cosmosClient) appendDelegationResponses(totalValidatorDelegations map[string]staking.DelegationResponse, validatorDelegations staking.DelegationResponses) {
	for _, d := range validatorDelegations {
		totalValidatorDelegations[d.GetDelegation().DelegatorAddress] = d
	}
}

func (c cosmosClient) ValidatorDelegations() (map[string]staking.DelegationResponse, error) {
	validatorDelegations := make(map[string]staking.DelegationResponse)
	var wg sync.WaitGroup

	// initial fetch to get total data
	res, err := c.validatorQueryClient.validatorDelegation(0, stride)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	c.appendDelegationResponses(validatorDelegations, res.GetDelegationResponses())

	total := res.GetPagination().GetTotal()
	fmt.Printf("total: %d\n", total)

	iterate := int(math.Ceil(float64(total / stride)))
	fmt.Printf("iterate: %d\n", iterate)

	wg.Add(iterate)

	// create iterate-sized buffered channel
	ch := make(chan staking.DelegationResponses, iterate)

	for i := 0; i < iterate; i++ {
		offset := i + 1
		go func() {
			defer wg.Done()
			vd, err := c.validatorQueryClient.validatorDelegation(uint64(offset*stride), stride)

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
		c.appendDelegationResponses(validatorDelegations, vd)
	}

	return validatorDelegations, err
}

func (c cosmosClient) ValidatorIncome() (*big.Int, error) {
	sdr, sdrErr := c.validatorQueryClient.selfDelegationReward()

	if sdrErr != nil {
		fmt.Println(sdrErr)
		return nil, sdrErr
	}

	cm, cmErr := c.validatorQueryClient.commission()

	if cmErr != nil {
		fmt.Println(cmErr)
		return nil, cmErr
	}

	commission := cm.GetCommission()
	income := commission.GetCommission()[0].Add(sdr.GetRewards()[0])

	return income.Amount.BigInt(), nil
}

func (c cosmosClient) GrantRewards() (map[string]*distribution.QueryDelegationRewardsResponse, error) {
	return c.grantQueryClient.rewards()
}

func (v validatorQueryClient) selfDelegationReward() (*distribution.QueryDelegationRewardsResponse, error) {
	res, err := v.distributionQueryClient.DelegationRewards(
		context.Background(),
		&distribution.QueryDelegationRewardsRequest{
			DelegatorAddress: v.validatorAddr,
			ValidatorAddress: v.validatorOperatorAddr,
		})
	return res, err
}

func (v validatorQueryClient) commission() (*distribution.QueryValidatorCommissionResponse, error) {
	res, err := v.distributionQueryClient.ValidatorCommission(
		context.Background(),
		&distribution.QueryValidatorCommissionRequest{
			ValidatorAddress: v.validatorOperatorAddr,
		})

	return res, err
}

// reward of grant wallet address delegated to given validator
func (g grantQueryClient) rewards() (map[string]*distribution.QueryDelegationRewardsResponse, error) {
	var wg sync.WaitGroup
	wg.Add(len(g.grantWalletAddr))

	delegationRewards := make(map[string]*distribution.QueryDelegationRewardsResponse)
	ch := make(chan *grantCommission, len(g.grantWalletAddr))

	for _, da := range g.grantWalletAddr {
		da := da
		go func() {
			defer wg.Done()
			fmt.Println(da)
			res, err := g.distributionClient.DelegationRewards(context.Background(), &distribution.QueryDelegationRewardsRequest{
				DelegatorAddress: da,
				ValidatorAddress: g.validatorOperatorAddr,
			})

			gc := grantCommission{
				da, res,
			}

			ch <- &gc

			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()
	close(ch)

	for r := range ch {
		delegationRewards[r.delegatorAddr] = r.commission
	}

	return delegationRewards, nil

}
