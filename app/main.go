package main

import (
	"fmt"
	"validator-dashboard/app/config"
	"validator-dashboard/app/server"
	"validator-dashboard/app/services/worker"
)

func main() {
	c := config.GetConfig()
	worker.RegisterCron(c.App.Cron)
	fmt.Println("Starting server...")

	app := server.NewApp()
	app.Run(fmt.Sprintf("%s:%s", c.App.Host, c.App.Port))
}
