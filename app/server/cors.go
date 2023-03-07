package server

import (
	"validator-dashboard/app/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCors(app *gin.Engine) {
	defaultConfig := cors.DefaultConfig()

	config := config.GetConfig()
	defaultConfig.AllowOrigins = config.App.AllowOrigins

	app.Use(cors.New(defaultConfig))
}
