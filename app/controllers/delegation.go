package controllers

import (
	"github.com/gin-gonic/gin"
)

func delegation(c *gin.Context) {
	chain := c.Param("chain")
	_ = chain
}

func AddDelegationRouters(r *gin.RouterGroup) {
	d := r.Group("/delegations")
	{
		d.GET("/:chain", delegation)
	}
}
