package app

import (
	"github.com/rs/zerolog/log"
	database "validator-dashboard/app/db"
)

func IncomeHistories(chain string) []database.IncomeHistory {
	db := database.GetDB()

	var incomeHistory []database.IncomeHistory

	query := `
		SELECT i.*, t.price, i.reward * price as reward_value, commission * price as commission_value
		FROM income_history i JOIN token_price t 
			ON 
				i.chain = t.chain 
				AND 
				DATE(i.create_dt) = DATE(t.create_dt) + INTERVAL '1 days' 
		WHERE i.chain = $1
		ORDER BY i.create_dt ASC
	`

	err := db.Select(&incomeHistory, query, chain)
	if err != nil {
		log.Err(err).Msg("failed to get income history")
	}

	return incomeHistory
}
