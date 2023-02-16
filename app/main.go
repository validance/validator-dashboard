package main

import (
	"fmt"
	"validator-dashboard/app/db"
	"validator-dashboard/app/services/client"
)

func main() {
	osmosisClient, err := client.NewCosmosClient("localhost:9090", "uatom", "cosmosvaloper1v78emy9d2xe3tj974l7tmn2whca2nh9zp7s0u9", "cosmos1v78emy9d2xe3tj974l7tmn2whca2nh9zy2y6sk")

	if err != nil {
		panic("failed to initialize client")
	}

	//delegations, err := osmosisClient.ValidatorDelegations()
	//fmt.Println(delegations["osmo18m4wkxw865cmxu7wv43pk9wgssw022kj732ee9"])

	validatorIncome, err := osmosisClient.ValidatorIncome()
	fmt.Println(validatorIncome)

	//rewards, err := osmosisClient.GrantRewards()
	//fmt.Println(rewards)
	db.New()

}
