package client

import (
	"validator-dashboard/app/models"
)

type aptosClient struct {
}

func (a aptosClient) ValidatorDelegations() (map[string]*models.Delegation, error) {
	return nil, nil
}
func (a aptosClient) ValidatorIncome() (*models.ValidatorIncome, error) {
	return nil, nil
}
func (a aptosClient) AddGrantAddresses([]string) {

}
func (a aptosClient) GrantRewards() (map[string]*models.Reward, error) {
	return nil, nil
}

func NewAptosClient() Client {
	return aptosClient{}
}
