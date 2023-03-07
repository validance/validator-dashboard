package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"validator-dashboard/app/services/client"
)

func test(c *gin.Context) {
	polygonClient, err := client.InitializePolygon()
	if err != nil {
		log.Err(err).Msg("failed to initialize polygon client")
	}
	polygonClient.ValidatorDelegations()
	polygonClient.ValidatorIncome()
	polygonClient.GrantRewards()

	c.JSON(200, gin.H{
		"message": "pong",
	})

}

func AddTestRouters(r *gin.RouterGroup) {
	d := r.Group("/polygon")
	{
		d.GET("/test", test)
	}
}
