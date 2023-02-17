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
		if client != nil {
			fmt.Printf("Initializing %s client\n", chain)
			cosmosChains = append(cosmosChains, client)
		}
	}
	return cosmosChains, nil
}

// TODO: initialize aptos client with config
func initializeAptos() (Client, error) {
	aptosConfig := config.GetConfig().Aptos
	_ = aptosConfig
	_ = NewAptosClient()
	// TODO: return nil when failed to initialize aptos client
	return nil, nil
}

// TODO: initialize polygon client with config
func initializePolygon() (Client, error) {
	polygonConfig := config.GetConfig().Polygon
	_ = polygonConfig
	_ = NewPolygonClient()
	// TODO: return nil when failed to initialize aptos client

	return nil, nil
}

func Initialize() ([]Client, error) {
	var clients []Client

	cosmosClients, cosmosErr := initializeCosmos()
	if cosmosErr != nil {
		fmt.Println(cosmosErr)
	}

	aptosClient, aptosErr := initializeAptos()
	if aptosErr != nil {
		fmt.Println(aptosErr)
	}

	polygonClient, polygonErr := initializePolygon()
	if polygonErr != nil {
		fmt.Println(polygonErr)
	}

	for _, c := range cosmosClients {
		if c != nil {
			clients = append(clients, c)
		}
	}

	if aptosClient != nil {
		clients = append(clients, aptosClient)
	}

	if polygonClient != nil {
		clients = append(clients, polygonClient)
	}

	return clients, nil
}
