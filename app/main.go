package main

import (
	"fmt"
	"validator-dashboard/app/config"
)

func main() {
	c := config.GetConfig()
	fmt.Println(c.Cosmos)

}
