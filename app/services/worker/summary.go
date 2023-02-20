package worker

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	database "validator-dashboard/app/db"
)

type SummaryWorker struct {
	db                 *sqlx.DB
	delegationChanged  []database.DelegationChanged
	newDelegators      []database.DelegationHistory
	leftDelegators     []database.DelegationHistory
	returnedDelegators []database.DelegationHistory
}

func NewSummaryWorker(dt *DelegationTask) *SummaryWorker {
	return &SummaryWorker{
		dt.db,
		dt.delegationChanged,
		dt.newDelegators,
		dt.leftDelegators,
		dt.returnedDelegators,
	}
}

func (s SummaryWorker) getAddressStatus(address string) (*database.AddressStatus, error) {
	var as []database.AddressStatus

	query := `
		SELECT *
		FROM address_status
		WHERE address = $1
	`

	err := s.db.Select(&as, query, address)
	if err != nil {
		log.Err(err).Msg("failed to get address status")
		return nil, err
	}

	if len(as) <= 0 {
		return nil, err
	}

	return &as[0], nil
}

func (s SummaryWorker) RunSummaryWorker() {
	for _, d := range s.delegationChanged {
		status, err := s.getAddressStatus(d.Address)
		if err != nil {
			continue
		}
		fmt.Println(status.Status)
	}
}
