package client

import (
	"validator-dashboard/app/models"
)

type polygonClient struct {
}

func (p polygonClient) ValidatorDelegations() (map[string]*models.Delegation, error) {
	return nil, nil
}
func (p polygonClient) ValidatorIncome() (*models.ValidatorIncome, error) {
	return nil, nil
}
func (p polygonClient) GrantRewards() (map[string]*models.GrantReward, error) {
	return nil, nil
}

func NewPolygonClient() Client {
	return aptosClient{}
}
