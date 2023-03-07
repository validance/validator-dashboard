package controllers

import (
	"net/http"
	"validator-dashboard/app/services/app"

	"github.com/gin-gonic/gin"
)

func address(c *gin.Context) {
	chain := c.Param("chain")
	addresses := app.AddressStatuses(chain)

	c.JSON(http.StatusOK, addresses)
}

func AddAddressRouters(r *gin.RouterGroup) {
	d := r.Group("/addresses")
	{
		d.GET("/:chain", address)
	}
}
