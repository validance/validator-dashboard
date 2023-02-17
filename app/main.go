package main

import (
	"github.com/rs/zerolog/log"
	"validator-dashboard/app/services/worker"
)

func main() {
	dbErr := worker.RunDbTask()
	if dbErr != nil {
		log.Err(dbErr)
	}
}
