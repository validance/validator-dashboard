package main

import (
	"validator-dashboard/app/server"
	"validator-dashboard/app/services/worker"
)

func main() {
	worker.Cron()

	app := server.NewApp()
	app.Run(":8000")
}
