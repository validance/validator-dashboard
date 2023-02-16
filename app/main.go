package main

import (
	"fmt"
	"validator-dashboard/app/db"
	"validator-dashboard/app/services/client"
)

func main() {
	grantAddrs := []string{"osmo1546y5dd36llwtmyl3pfrxtq7tzys5853543h7r", "osmo1tdcu25km6am5ya6armwne995zarjzzv6g5why0"}

	osmosisClient, err := client.NewCosmosClient("osmosis-grpc.polkachu.com:12590", "osmovaloper18m4wkxw865cmxu7wv43pk9wgssw022kjyxz6wz", "osmo18m4wkxw865cmxu7wv43pk9wgssw022kj732ee9", grantAddrs...)

	if err != nil {
		panic("failed to initialize client")
	}

	delegations, err := osmosisClient.ValidatorDelegations()
	fmt.Print(delegations["osmo18m4wkxw865cmxu7wv43pk9wgssw022kj732ee9"])

	//validatorIncome, err := osmosisClient.ValidatorIncome()
	//fmt.Println(validatorIncome)

	//rewards, err := osmosisClient.GrantRewards()
	//fmt.Println(rewards)

	database := db.New()

	fmt.Println(database)
}
