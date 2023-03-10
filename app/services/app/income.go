package app

import (
	database "validator-dashboard/app/db"

	"github.com/rs/zerolog/log"
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

func IncomeSummaryByDate(date string) database.IncomeSummary {
	db := database.GetDB()

	var incomeSummary database.IncomeSummary

	query := `
		SELECT 
			SUM(i.reward * t.price) as total_reward_value,
			SUM(i.commission * t.price) as total_commission_value,
			SUM(g.reward_sum * t.price) as total_grant_reward_value
		FROM income_history i
			JOIN token_price t
			ON 
				i.chain = t.chain
				AND
				DATE(i.create_dt) = DATE(t.create_dt) + INTERVAL '1 days'
			LEFT JOIN (
				SELECT SUM(reward) reward_sum, chain
				FROM grant_reward_history
				WHERE DATE(create_dt) = $1
				GROUP BY chain
			) g
			ON 
				i.chain = g.chain
		WHERE DATE(i.create_dt) = $1
		GROUP BY DATE(i.create_dt)
	`

	err := db.Get(&incomeSummary, query, date)
	if err != nil {
		log.Err(err).Msg("failed to get income summary")
	}

	return incomeSummary
}
