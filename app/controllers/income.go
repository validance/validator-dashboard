package controllers

import (
	"net/http"
	"validator-dashboard/app/services/app"

	"github.com/gin-gonic/gin"
)

func getIncomeHistories(c *gin.Context) {
	chain := c.Param("chain")

	incomeHistories := app.IncomeHistories(chain)
	c.JSON(http.StatusOK, incomeHistories)
}

func getIncomeSummary(c *gin.Context) {
	date := c.Param("date")

	incomeSummary := app.IncomeSummaryByDate(date)
	c.JSON(http.StatusOK, incomeSummary)
}

func AddIncomeRouters(r *gin.RouterGroup) {
	d := r.Group("/validator-income")
	{
		d.GET("/chains/:chain", getIncomeHistories)
	}
	{
		d.GET("/summaries/:date", getIncomeSummary)
	}
}
