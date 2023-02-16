package client

import "math/big"

type polygonClient struct {
}

func (a polygonClient) ValidatorDelegations() (map[string]*big.Int, error) {
	return nil, nil
}
func (a polygonClient) ValidatorIncome() (*big.Int, error) {
	return nil, nil
}
func (a polygonClient) AddGrantAddresses([]string) {

}
func (a polygonClient) GrantRewards() (map[string]*big.Int, error) {
	return nil, nil
}

func NewPolygonClient() Client {
	return aptosClient{}
}
