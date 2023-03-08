package app

import (
	"database/sql"
	database "validator-dashboard/app/db"

	"github.com/rs/zerolog/log"
)

func GetAddressStatuses(chain string) []database.AddressStatus {
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

func UpdateAddressLabel(address string, label string) sql.Result {
	db := database.GetDB()

	query := `
		UPDATE address_status
		SET
			label = $2
		WHERE address = $1
	`

	res, err := db.Exec(query, address, label)
	if err != nil {
		log.Err(err).Msg("failed to update address label")
	}

	return res
}
