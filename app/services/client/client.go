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
	GrantRewards() (map[string]*models.GrantReward, error)
}

func InitializeCosmos() ([]Client, error) {
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
func InitializeAptos() (Client, error) {
	aptosConfig := config.GetConfig().Aptos
	_ = aptosConfig
	_ = NewAptosClient()
	// TODO: return nil when failed to initialize aptos client
	return nil, nil
}

// TODO: initialize polygon client with config
func InitializePolygon() (Client, error) {
	polygonConfig := config.GetConfig().Polygon
	client, err := NewPolygonClient(polygonConfig.EndpointUrl, polygonConfig.Denom, polygonConfig.Exponent, polygonConfig.OwnerAddr, polygonConfig.ValidatorIndex)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Initialize() []Client {
	var clients []Client

	cosmosClients, cosmosErr := InitializeCosmos()
	if cosmosErr != nil {
		log.Err(cosmosErr).Msg("failed to initialize cosmos clients")
	}

	aptosClient, aptosErr := InitializeAptos()
	if aptosErr != nil {
		log.Err(aptosErr).Msg("failed to initialize aptos client")
	}

	polygonClient, polygonErr := InitializePolygon()
	if polygonErr != nil {
		log.Err(polygonErr).Msg("failed to initialize polygon client")
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

	return clients
}
