package controllers

import (
	"net/http"
	"validator-dashboard/app/services/app"

	"github.com/gin-gonic/gin"
)

func grantReward(c *gin.Context) {
	chain := c.Param("chain")

	grantRewardHistories := app.GrantRewardHistories(chain)
	c.JSON(http.StatusOK, grantRewardHistories)
}

func AddGrantRewardRouters(r *gin.RouterGroup) {
	d := r.Group("/grant-rewards")
	{
		d.GET("/:chain", grantReward)
	}
}
