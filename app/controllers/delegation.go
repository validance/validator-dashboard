package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services"
)

func delegation(c *gin.Context) {
	chain := c.Param("chain")
	summary, check := services.GetCache("delegation_summary")
	if !check {
		log.Error().Msg("summary not found")
	}

	sum, err := summary.(map[string]*models.DelegationSummary)
	if !err {
		c.String(http.StatusNotFound, "summary not found")
	} else {
		c.JSON(http.StatusOK, sum[chain])
	}
}

func AddDelegationRouters(r *gin.RouterGroup) {
	d := r.Group("/delegations")
	{
		d.GET("/:chain", delegation)
	}
}
