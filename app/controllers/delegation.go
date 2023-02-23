package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"validator-dashboard/app/services/app"
)

func delegation(c *gin.Context) {
	chain := c.Param("chain")
	histories := app.DelegationSummaryHistories(chain)

	c.JSON(http.StatusOK, histories)
}

func AddDelegationRouters(r *gin.RouterGroup) {
	d := r.Group("/delegations")
	{
		d.GET("/:chain", delegation)
	}
}
