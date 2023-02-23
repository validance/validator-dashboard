package app

import (
	"github.com/rs/zerolog/log"
	database "validator-dashboard/app/db"
)

func DelegationSummaryHistories(chain string) []database.DelegationSummaryHistory {
	db := database.GetDB()

	var delegationSummaryHistories []database.DelegationSummaryHistory

	query := `
		SELECT 
			d.*, 
			t.price as price,
			
			d.yesterday_delegation_amount_total * t.price as yesterday_delegation_value_total,
			d.yesterday_delegation_amount_b2b * t.price as yesterday_delegation_value_b2b,
			d.yesterday_delegation_amount_b2c * t.price as yesterday_delegation_value_b2c,
			d.yesterday_delegation_amount_unknown * t.price as yesterday_delegation_value_unknown,
			
			d.today_existing_increased_delegation_amount_total * t.price as today_existing_increased_delegation_value_total,
			d.today_existing_increased_delegation_amount_b2b * t.price as today_existing_increased_delegation_value_b2b,
			d.today_existing_increased_delegation_amount_b2c * t.price as today_existing_increased_delegation_value_b2c,
			d.today_existing_increased_delegation_amount_unknown * t.price as today_existing_increased_delegation_value_unknown,
			
			d.today_new_increased_delegation_amount_total * t.price as today_new_increased_delegation_value_total,
			d.today_new_increased_delegation_amount_b2b  * t.price as today_new_increased_delegation_value_b2b,
			d.today_new_increased_delegation_amount_b2c * t.price as today_new_increased_delegation_value_b2c,
			d.today_new_increased_delegation_amount_unknown * t.price as today_new_increased_delegation_value_unknown,
			
			d.today_return_increased_delegation_amount_total * t.price as today_return_increased_delegation_value_total,
			d.today_return_increased_delegation_amount_b2b * t.price as today_return_increased_delegation_value_b2b,
			d.today_return_increased_delegation_amount_b2c * t.price as today_return_increased_delegation_value_b2c,
			d.today_return_increased_delegation_amount_unknown * t.price as today_return_increased_delegation_value_unknown,
			
			d.today_existing_decreased_delegation_amount_total * t.price as today_existing_decreased_delegation_value_total,
			d.today_existing_decreased_delegation_amount_b2b * t.price as today_existing_decreased_delegation_value_b2b,
			d.today_existing_decreased_delegation_amount_b2c * t.price as today_existing_decreased_delegation_value_b2c,
			d.today_existing_decreased_delegation_amount_unknown * t.price as today_existing_decreased_delegation_value_unknown,
			
			d.today_left_decreased_delegation_amount_total * t.price as today_left_decreased_delegation_value_total,
			d.today_left_decreased_delegation_amount_b2b * t.price as today_left_decreased_delegation_value_b2b,
			d.today_left_decreased_delegation_amount_b2c * t.price as today_left_decreased_delegation_value_b2c,
			d.today_left_decreased_delegation_amount_unknown * t.price as today_left_decreased_delegation_value_unknown
			
		FROM delegation_summary d
			JOIN token_price t
			ON 
				d.chain = t.chain
				AND
				DATE(d.create_dt) = DATE(t.create_dt) + INTERVAL '1 days'
		WHERE d.chain = $1
	`

	err := db.Select(&delegationSummaryHistories, query, chain)
	if err != nil {
		log.Err(err).Msg("failed to get delegation summary history")
	}

	return delegationSummaryHistories
}
