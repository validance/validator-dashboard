package client

import (
	"fmt"
	"math/big"
	"validator-dashboard/app/config"
)

type Client interface {
	ValidatorDelegations() (map[string]*big.Int, error)
	ValidatorIncome() (*big.Int, error)
	// AddGrantAddresses add address delegated to given validator
	AddGrantAddresses([]string)
	// GrantRewards get reward per each grant address
	GrantRewards() (map[string]*big.Int, error)
}

func initializeCosmos() ([]Client, error) {
	cosmosConfig := config.GetConfig().Cosmos
	cosmosChains := make([]Client, len(cosmosConfig))

	for chain, info := range cosmosConfig {
		client, err := NewCosmosClient(info.GrpcUrl, info.Denom, info.ValidatorOperatorAddr, info.ValidatorAddr)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Initializing %s client\n", chain)
		cosmosChains = append(cosmosChains, client)
	}
	return cosmosChains, nil
}

// TODO: initialize aptos client with config
func initializeAptos() (Client, error) {
	aptosConfig := config.GetConfig().Aptos
	_ = aptosConfig
	return NewAptosClient(), nil
}

// TODO: initialize polygon client with config
func initializePolygon() (Client, error) {
	polygonConfig := config.GetConfig().Polygon
	_ = polygonConfig
	return NewPolygonClient(), nil
}

func Initialize() ([]Client, error) {
	clients := make([]Client, 1)

	cosmosClients, cosmosErr := initializeCosmos()
	if cosmosErr != nil {
		return nil, cosmosErr
	}

	aptosClient, aptosErr := initializeAptos()
	if cosmosErr != nil {
		return nil, aptosErr
	}

	polygonClient, polygonErr := initializePolygon()
	if cosmosErr != nil {
		return nil, polygonErr
	}

	for _, c := range cosmosClients {
		clients = append(clients, c)
	}

	clients = append(clients, aptosClient)
	clients = append(clients, polygonClient)

	return clients, nil
}
