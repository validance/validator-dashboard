package controllers

import (
	"net/http"
	"validator-dashboard/app/services/app"

	"github.com/gin-gonic/gin"
)

func getDelegationHistories(c *gin.Context) {
	chain := c.Param("chain")
	histories := app.DelegationSummaryHistoriesByChain(chain)

	c.JSON(http.StatusOK, histories)
}

func getDelegationSummary(c *gin.Context) {
	date := c.Param("date")
	summary := app.DelegationSummaryByDate(date)

	c.JSON(http.StatusOK, summary)
}

func AddDelegationRouters(r *gin.RouterGroup) {
	d := r.Group("/delegations")
	{
		d.GET("/chains/:chain", getDelegationHistories)
	}
	{
		d.GET("/summaries/:date", getDelegationSummary)
	}
}
