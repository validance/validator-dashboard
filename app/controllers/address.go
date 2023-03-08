package controllers

import (
	"net/http"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services/app"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func getAddressesByChain(c *gin.Context) {
	chain := c.Param("chain")
	addresses := app.GetAddressStatuses(chain)

	c.JSON(http.StatusOK, addresses)
}

func patchAddress(c *gin.Context) {
	address := c.Param("address")

	var body models.PatchAddressBody
	c.BindWith(&body, binding.JSON)

	app.UpdateAddressLabel(address, body.Label)

	c.Status(http.StatusOK)
}

func AddAddressRouters(r *gin.RouterGroup) {
	d := r.Group("/addresses")
	{
		d.PATCH("/:address", patchAddress)
	}
	{
		d.GET("/chains/:chain", getAddressesByChain)
	}
}
