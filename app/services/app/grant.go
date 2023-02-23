package app

import (
	database "validator-dashboard/app/db"

	"github.com/rs/zerolog/log"
)

func GrantRewardHistories(chain string) []database.GrantRewardHistory {
	db := database.GetDB()

	var grantRewardHistories []database.GrantRewardHistory

	query := `
	SELECT g.*, g.reward * p.price as reward_value
	FROM grant_reward_history g JOIN token_price p
		ON
			g.chain = p.chain
			AND
			DATE(g.create_dt) = DATE(p.create_dt) + INTERVAL '1 days'
	WHERE g.chain = $1
	ORDER BY g.create_dt ASC
	`

	err := db.Select(&grantRewardHistories, query, chain)
	if err != nil {
		log.Err(err).Msg("failed to get grant reward history")
	}

	return grantRewardHistories
}
