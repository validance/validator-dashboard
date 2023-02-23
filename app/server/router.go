package server

import (
	"github.com/gin-gonic/gin"
	"validator-dashboard/app/controllers"
)

func AddRouters(app *gin.Engine) {
	v1 := app.Group("/api/v1")

	controllers.AddDelegationRouters(v1)
	controllers.AddIncomeRouters(v1)
}
