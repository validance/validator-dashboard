package client

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"math/big"
	"sync"
	"validator-dashboard/app/models"
)

const stride = 1000

// coefficient to divide dec coin to normalize into denomination
var coin_c = big.NewInt(int64(math.Pow10(18)))

type grantCommission struct {
	delegatorAddr string
	commission    *distribution.QueryDelegationRewardsResponse
}

type validatorQuerier interface {
	validatorDelegations(offset uint64, limit uint64) (*staking.QueryValidatorDelegationsResponse, error)
	selfDelegationReward() (*distribution.QueryDelegationRewardsResponse, error)
	commission() (*distribution.QueryValidatorCommissionResponse, error)
	getValidatorAddr() string
	getOperatorAddr() string
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
	chain                string
	denom                string
	exponent             *big.Float
	url                  string
	validatorQueryClient validatorQuerier
	grantQueryClient     grantQuerier
}

// NewCosmosClient create query client for cosmos app-chains
func NewCosmosClient(url string, chain string, denom string, exponent int, validatorOperatorAddr string, validatorAddr string, grantWalletAddr ...string) (Client, error) {
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
		chain,
		denom,
		new(big.Float).SetInt(big.NewInt(int64(math.Pow10(exponent)))),
		url,
		vqc,
		gqc,
	}, nil
}

func (v validatorQueryClient) validatorDelegations(offset uint64, limit uint64) (*staking.QueryValidatorDelegationsResponse, error) {
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

// divide by coin coefficient and push the value into map
func (c cosmosClient) appendDelegationResponses(totalValidatorDelegations map[string]*models.Delegation, validatorDelegations staking.DelegationResponses) {
	for _, d := range validatorDelegations {
		delegationAmount := d.GetDelegation().GetShares().BigInt()
		delegationAmountBigF := BigIntToFloat(delegationAmount.Div(delegationAmount, coin_c))
		delegationAmountBigF = delegationAmountBigF.Quo(delegationAmountBigF, c.exponent)
		delegationAmountF, _ := delegationAmountBigF.Float64()

		if FilterLowAmount(delegationAmountBigF) {
			continue
		}

		delegation := &models.Delegation{
			Address:   d.GetDelegation().DelegatorAddress,
			Validator: d.GetDelegation().ValidatorAddress,
			Chain:     c.chain,
			Amount:    delegationAmountF,
		}

		totalValidatorDelegations[d.GetDelegation().DelegatorAddress] = delegation
	}
}

func (c cosmosClient) ValidatorDelegations() (map[string]*models.Delegation, error) {
	validatorDelegations := make(map[string]*models.Delegation)
	var wg sync.WaitGroup

	// initial fetch to get total data
	res, err := c.validatorQueryClient.validatorDelegations(0, stride)

	if err != nil {
		log.Err(err).Msg("failed to get cosmos validator delegations")
		return nil, err
	}

	c.appendDelegationResponses(validatorDelegations, res.GetDelegationResponses())

	total := res.GetPagination().GetTotal()
	iterate := int(math.Ceil(float64(total / stride)))

	wg.Add(iterate)

	// create iterate-sized buffered channel
	ch := make(chan staking.DelegationResponses, iterate)

	for i := 0; i < iterate; i++ {
		offset := i + 1
		go func() {
			defer wg.Done()
			vd, err := c.validatorQueryClient.validatorDelegations(uint64(offset*stride), stride)

			if err != nil {
				log.Err(err).Msg("failed to get cosmos validator delegations")
			}
			ch <- vd.GetDelegationResponses()
		}()
	}

	// wait until every goroutine finish pushing data into channel
	wg.Wait()

	// close buffered channel
	close(ch)

	// pop data from buffered channel
	for vds := range ch {
		c.appendDelegationResponses(validatorDelegations, vds)
	}

	return validatorDelegations, err
}

func (c cosmosClient) ValidatorIncome() (*models.ValidatorIncome, error) {
	sdr, sdrErr := c.validatorQueryClient.selfDelegationReward()

	if sdrErr != nil {
		return nil, sdrErr
	}

	cm, cmErr := c.validatorQueryClient.commission()

	if cmErr != nil {
		return nil, cmErr
	}

	reward := sdr.GetRewards().AmountOf(c.denom).BigInt()
	rewardValBigF := BigIntToFloat(reward.Div(reward, coin_c))
	rewardValBigF = rewardValBigF.Quo(rewardValBigF, c.exponent)
	rewardValF, _ := rewardValBigF.Float64()

	commission := cm.GetCommission()
	commissionVal := commission.GetCommission().AmountOf(c.denom).BigInt()
	commissionValBigF := BigIntToFloat(commissionVal.Div(commissionVal, coin_c))
	commissionValBigF = commissionValBigF.Quo(commissionValBigF, c.exponent)
	commissionValF, _ := commissionValBigF.Float64()

	validatorIncome := &models.ValidatorIncome{
		Chain:      c.chain,
		Validator:  c.validatorQueryClient.getValidatorAddr(),
		Reward:     rewardValF,
		Commission: commissionValF,
	}

	return validatorIncome, nil
}

func (c cosmosClient) GrantRewards() (map[string]*models.GrantReward, error) {
	res := make(map[string]*models.GrantReward)
	rewards, err := c.grantQueryClient.rewards()

	if err != nil {
		return nil, err
	}

	for delegatorAddr, reward := range rewards {
		if reward != nil {
			rewardVal := &models.GrantReward{
				Chain:     c.chain,
				Validator: c.validatorQueryClient.getOperatorAddr(),
				Reward:    0,
			}
			r := reward.GetRewards().AmountOf(c.denom).BigInt()
			rf := BigIntToFloat(r.Div(r, coin_c))
			rf = rf.Quo(rf, c.exponent)
			f, _ := rf.Float64()
			rewardVal.Reward = f

			res[delegatorAddr] = rewardVal
		}
	}

	return res, nil
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

func (v validatorQueryClient) getValidatorAddr() string {
	return v.validatorAddr
}

func (v validatorQueryClient) getOperatorAddr() string {
	return v.validatorOperatorAddr
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
			res, err := g.distributionClient.DelegationRewards(context.Background(), &distribution.QueryDelegationRewardsRequest{
				DelegatorAddress: da,
				ValidatorAddress: g.validatorOperatorAddr,
			})
			gc := grantCommission{da, res}
			ch <- &gc

			if err != nil {
				log.Err(err).Msg("failed to get cosmos grant rewards")
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
