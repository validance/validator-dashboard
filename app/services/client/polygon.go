package client

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"validator-dashboard/app/models"
)

type polygonClient struct {
	denom                 string
	exponent              *big.Float
	baseUrl               string
	valdatorInfoUrl       string
	delegatorUrl          string
	commissionedRewardUrl string
	ownerAddr             string
	validatorIndex        int
	chain                 string
	polygonDelegatorInfo  PolygonDelegatorInfo
	validatorInfo         ValidatorInfo
	commissionInfo        CommissionedReward
}
type PolygonDelegatorInfo struct {
	Summary struct {
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		SortBy    string `json:"sortBy"`
		Direction string `json:"direction"`
		Total     int    `json:"total"`
		Size      int    `json:"size"`
	} `json:"summary"`
	Success bool        `json:"success"`
	Status  string      `json:"status"`
	Result  []delegator `json:"result"`
}
type delegator struct {
	BondedValidator   int     `json:"bondedValidator"`
	Stake             float64 `json:"stake"`
	Address           string  `json:"address"`
	ClaimedReward     int64   `json:"claimedReward"`
	Shares            string  `json:"shares"`
	DeactivationEpoch string  `json:"deactivationEpoch"`
}

type ValidatorInfo struct {
	Status  string `json:"status"`
	Success bool   `json:"success"`
	Result  struct {
		Id                                 int     `json:"id"`
		Name                               string  `json:"name"`
		Description                        string  `json:"description"`
		Url                                string  `json:"url"`
		LogoUrl                            string  `json:"logoUrl"`
		Owner                              string  `json:"owner"`
		Signer                             string  `json:"signer"`
		CommissionPercent                  string  `json:"commissionPercent"`
		SignerPublicKey                    string  `json:"signerPublicKey"`
		SelfStake                          string  `json:"selfStake"`
		DelegatedStake                     string  `json:"delegatedStake"`
		IsInAuction                        bool    `json:"isInAuction"`
		AuctionAmount                      string  `json:"auctionAmount"`
		ClaimedReward                      string  `json:"claimedReward"`
		ActivationEpoch                    string  `json:"activationEpoch"`
		TotalStaked                        string  `json:"totalStaked"`
		DeactivationEpoch                  string  `json:"deactivationEpoch"`
		JailEndEpoch                       string  `json:"jailEndEpoch"`
		Status                             string  `json:"status"`
		ContractAddress                    string  `json:"contractAddress"`
		UptimePercent                      float64 `json:"uptimePercent"`
		DelegationEnabled                  bool    `json:"delegationEnabled"`
		DelegatorCount                     int     `json:"delegatorCount"`
		DelegatorUnclaimedRewards          string  `json:"delegatorUnclaimedRewards"`
		ValidatorUnclaimedRewards          string  `json:"validatorUnclaimedRewards"`
		DelegatorClaimedRewards            float64 `json:"delegatorClaimedRewards"`
		CheckpointsMissed                  int     `json:"checkpointsMissed"`
		CheckpointsSigned                  int     `json:"checkpointsSigned"`
		MissedLatestCheckpointcount        int     `json:"missedLatestCheckpointcount"`
		PerformanceIndex                   float64 `json:"performanceIndex"`
		LastConsideredEngagementCheckpoint int     `json:"lastConsideredEngagementCheckpoint"`
		CurrentState                       string  `json:"currentState"`
	} `json:"result"`
}
type CommissionedReward struct {
	Status  string `json:"status"`
	Success bool   `json:"success"`
	Result  struct {
		TotalCommissionedRewards float64 `json:"totalCommissionedRewards"`
	} `json:"result"`
}

// divide by coin coefficient and push the value into map
func (p polygonClient) appendDelegationResponses(totalValidatorDelegations map[string]*models.Delegation, validatorDelegations []delegator) {
	for _, d := range validatorDelegations {
		delegationAmount := d.Stake
		//delegationAmountBigF := delegationAmount.Quo(delegationAmount, p.exponent)
		//delegationAmountF, _ := delegationAmountBigF.Float64()

		delegation := &models.Delegation{
			Address:   d.Address,
			Validator: p.ownerAddr,
			Chain:     p.chain,
			Amount:    delegationAmount / math.Pow10(18),
		}

		totalValidatorDelegations[d.Address] = delegation
	}
}

func (p polygonClient) ValidatorDelegations() (map[string]*models.Delegation, error) {
	validatorDelegations := make(map[string]*models.Delegation)

	validatorInfo := p.validatorInfo
	delegatorsInfo := p.polygonDelegatorInfo

	//self stake info add
	selfStakeBigF := StringToFloat(validatorInfo.Result.SelfStake)
	selfStakeF, _ := selfStakeBigF.Float64()

	delegator := delegator{
		BondedValidator: validatorInfo.Result.Id,
		Stake:           selfStakeF,
		Address:         validatorInfo.Result.Owner,
	}

	delegators := delegatorsInfo.Result
	delegators = append(delegators, delegator)
	log.Info().Msgf("found %d delegators", len(delegators))
	p.appendDelegationResponses(validatorDelegations, delegators)
	// print validatorDelegations values
	//for _, v := range validatorDelegations {
	//	log.Info().Msgf("delegation: %s", v)
	//}

	return validatorDelegations, nil
}
func (p polygonClient) ValidatorIncome() (*models.ValidatorIncome, error) {

	validatorInfo := p.validatorInfo.Result
	commissionInfo := p.commissionInfo.Result

	rewardBigF := StringToFloat(validatorInfo.ValidatorUnclaimedRewards)
	rewardBigF = rewardBigF.Quo(rewardBigF, p.exponent)
	rewardF, _ := rewardBigF.Float64()

	//commissionBigF := StringToFloat(commissionInfo.TotalCommissionedRewards)
	//commissionBigF = commissionBigF.Quo(commissionBigF, p.exponent)
	//commissionF, _ := commissionBigF.Float64()

	validatorIncome := &models.ValidatorIncome{
		Validator:  p.ownerAddr,
		Chain:      p.chain,
		Reward:     rewardF / math.Pow10(18),
		Commission: commissionInfo.TotalCommissionedRewards / math.Pow10(18),
	}

	// print validatorIncome values
	//log.Info().Msgf("validatorIncome: %s", validatorIncome)

	return validatorIncome, nil
}
func (p polygonClient) GrantRewards() (map[string]*models.GrantReward, error) {
	res := make(map[string]*models.GrantReward)
	return res, nil
}

func NewPolygonClient(url string, denom string, exponent int, ownerAddr string, validatorIndex int) (Client, error) {
	client := polygonClient{
		denom:                 denom,
		exponent:              new(big.Float).SetInt(big.NewInt(int64(math.Pow10(exponent)))),
		baseUrl:               url,
		valdatorInfoUrl:       url + "validators/" + strconv.Itoa(validatorIndex),
		delegatorUrl:          url + "validators/" + strconv.Itoa(validatorIndex) + "/delegators?offset=0&limit=100",
		commissionedRewardUrl: url + "validators/" + strconv.Itoa(validatorIndex) + "/commissioned-reward",
		ownerAddr:             ownerAddr,
		validatorIndex:        validatorIndex,
		chain:                 "polygon",
	}

	// get delegator list
	log.Info().Msg("getting polygon delegator list. " + client.delegatorUrl)
	resp, err := http.Get(client.delegatorUrl)
	if err != nil {
		log.Err(err).Msg("failed to get delegator list")
		return nil, err
	}
	var delegatorsInfo PolygonDelegatorInfo
	err = json.NewDecoder(resp.Body).Decode(&delegatorsInfo)
	if err != nil {
		log.Err(err).Msg("failed to decode delegator list")
		return nil, err
	}
	defer resp.Body.Close()

	client.polygonDelegatorInfo = delegatorsInfo

	//get validator info
	log.Info().Msg("getting polygon validator info. " + client.valdatorInfoUrl)
	resp, err = http.Get(client.valdatorInfoUrl)
	if err != nil {
		log.Err(err).Msg("failed to get validator info")
		return nil, err
	}
	var validatorInfo ValidatorInfo
	err = json.NewDecoder(resp.Body).Decode(&validatorInfo)
	if err != nil {
		log.Err(err).Msg("failed to decode validator info")
		return nil, err
	}
	defer resp.Body.Close()

	client.validatorInfo = validatorInfo

	//get commission info
	log.Info().Msg("getting polygon commission info. " + client.commissionedRewardUrl)
	resp, err = http.Get(client.commissionedRewardUrl)
	if err != nil {
		log.Err(err).Msg("failed to get commission info")
		return nil, err
	}
	var commissionInfo CommissionedReward
	err = json.NewDecoder(resp.Body).Decode(&commissionInfo)
	if err != nil {
		log.Err(err).Msg("failed to decode commission info")
		return nil, err
	}
	defer resp.Body.Close()

	// print commissionInfo values
	log.Info().Msgf("commissionInfo: %s", commissionInfo)

	client.commissionInfo = commissionInfo

	return client, nil

}
