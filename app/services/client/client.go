package client

import (
	"fmt"
	"validator-dashboard/app/config"
	"validator-dashboard/app/models"
)

type Client interface {
	// ValidatorDelegations includes validator self-bonded tokens
	ValidatorDelegations() (map[string]*models.Delegation, error)
	ValidatorIncome() (*models.ValidatorIncome, error)
	// AddGrantAddresses add address delegated to given validator
	AddGrantAddresses([]string)
	// GrantRewards get reward per each grant address
	GrantRewards() (map[string]*models.Reward, error)
}

func initializeCosmos() ([]Client, error) {
	var cosmosChains []Client

	cosmosConfig := config.GetConfig().Cosmos

	for chain, info := range cosmosConfig {
		client, err := NewCosmosClient(info.GrpcUrl, chain, info.Denom, info.ValidatorOperatorAddr, info.ValidatorAddr)
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
	var clients []Client

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
