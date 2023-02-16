package main

import (
	"fmt"
	"sync"
	"validator-dashboard/app/services/client"
)

func main() {
	clients, err := client.Initialize()

	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(clients))

	for _, c := range clients {
		c := c
		go func() {
			defer wg.Done()
			vi, err := c.ValidatorDelegations()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%v\n", vi)
		}()
	}

	wg.Wait()
}
