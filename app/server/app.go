package server

import (
	"github.com/gin-gonic/gin"
)

func NewApp() *gin.Engine {
	app := gin.Default()

	SetupCors(app)
	AddRouters(app)

	return app
}
