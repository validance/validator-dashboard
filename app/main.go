package main

import (
	"github.com/rs/zerolog/log"
	"validator-dashboard/app/services/worker"
)

func main() {
	dbErr := worker.Run()
	if dbErr != nil {
		log.Err(dbErr)
	}
}
