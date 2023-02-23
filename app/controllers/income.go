package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"validator-dashboard/app/services/app"
)

func income(c *gin.Context) {
	chain := c.Param("chain")

	incomeHistories := app.IncomeHistories(chain)
	c.JSON(http.StatusOK, incomeHistories)
}

func AddIncomeRouters(r *gin.RouterGroup) {
	d := r.Group("/validator-income")
	{
		d.GET("/:chain", income)
	}
}
