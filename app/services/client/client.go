package client

import (
	"github.com/rs/zerolog/log"
	"validator-dashboard/app/config"
	"validator-dashboard/app/models"
)

type Client interface {
	// ValidatorDelegations includes validator self-bonded tokens
	ValidatorDelegations() (map[string]*models.Delegation, error)
	ValidatorIncome() (*models.ValidatorIncome, error)
	// GrantRewards get reward per each grant address
	GrantRewards() (map[string]*models.Reward, error)
}

func initializeCosmos() ([]Client, error) {
	var cosmosChains []Client

	cosmosConfig := config.GetConfig().Cosmos

	for chain, info := range cosmosConfig {
		client, err := NewCosmosClient(info.GrpcUrl, chain, info.Denom, info.Exponent, info.ValidatorOperatorAddr, info.ValidatorAddr, info.GrantAddrs...)
		if err != nil {
			return nil, err
		}
		if client != nil {
			log.Info().Msgf("Initializing %s client", chain)
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
		log.Err(cosmosErr)
	}

	aptosClient, aptosErr := initializeAptos()
	if aptosErr != nil {
		log.Err(aptosErr)
	}

	polygonClient, polygonErr := initializePolygon()
	if polygonErr != nil {
		log.Err(polygonErr)
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
