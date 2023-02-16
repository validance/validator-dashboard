package client

import "math/big"

type aptosClient struct {
}

func (a aptosClient) ValidatorDelegations() (map[string]*big.Int, error) {
	return nil, nil
}
func (a aptosClient) ValidatorIncome() (*big.Int, error) {
	return nil, nil
}
func (a aptosClient) AddGrantAddresses([]string) {

}
func (a aptosClient) GrantRewards() (map[string]*big.Int, error) {
	return nil, nil
}

func NewAptosClient() Client {
	return aptosClient{}
}
