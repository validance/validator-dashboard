package server

import "github.com/gin-gonic/gin"

func NewApp() *gin.Engine {
	app := gin.Default()
	AddRouters(app)

	return app
}
