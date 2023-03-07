package app

import (
	database "validator-dashboard/app/db"

	"github.com/rs/zerolog/log"
)

func AddressStatuses(chain string) []database.AddressStatus {
	db := database.GetDB()

	var addressStatus []database.AddressStatus

	query := `
		SELECT *
		FROM address_status
		WHERE chain = $1
		ORDER BY create_dt DESC
	`

	err := db.Select(&addressStatus, query, chain)
	if err != nil {
		log.Err(err).Msg("failed to get address status")
	}

	return addressStatus
}
