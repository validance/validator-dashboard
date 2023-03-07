package server

import (
	"validator-dashboard/app/controllers"

	"github.com/gin-gonic/gin"
)

func AddRouters(app *gin.Engine) {
	v1 := app.Group("/api/v1")

	controllers.AddDelegationRouters(v1)
	controllers.AddIncomeRouters(v1)
	controllers.AddGrantRewardRouters(v1)

	test := app.Group("/test")
	controllers.AddTestRouters(test)
}
