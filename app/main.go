package main

import (
	"fmt"
	"validator-dashboard/app/services/client"
)

func main() {
	osmosisClient, err := client.NewCosmosClient("osmosis-grpc.polkachu.com:12590", "osmovaloper18m4wkxw865cmxu7wv43pk9wgssw022kjyxz6wz")

	if err != nil {
		panic("failed to initialize client")
	}

	delegations, err := osmosisClient.ValidatorDelegations()

	fmt.Print(len(delegations))
}
