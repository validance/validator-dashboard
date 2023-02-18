package worker

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type delegationTask struct {
	db *sqlx.DB
}

func newDelegationTask(db *sqlx.DB) *delegationTask {
	return &delegationTask{db}
}

func (d delegationTask) getManagedChains() []string {
	var chains []string

	queryErr := d.db.Select(
		&chains,
		`
			SELECT DISTINCT chain
			FROM delegation_history
		`,
	)

	if queryErr != nil {
		log.Err(queryErr)
	}

	return chains
}

func (d delegationTask) getNewDelegators() []string {
	var newDelegators []string

	newDelegatorQuery := `
		SELECT address
		FROM delegation_history
		WHERE NOT EXISTS (
			SELECT *
			FROM address_status
			WHERE 
			    address_status.address = delegation_history.address
		)
	`

	err := d.db.Select(&newDelegators, newDelegatorQuery)
	if err != nil {
		log.Err(err)
	}

	return newDelegators
}

func (d delegationTask) getLeftDelegators() []string {
	var leftDelegators []string

	leftDelegatorsQuery := `
		SELECT address 
		FROM address_status
		WHERE NOT EXISTS (
			SELECT *
			FROM delegation_history
			WHERE address_status.address = delegation_history.address
		)
	`

	err := d.db.Select(&leftDelegators, leftDelegatorsQuery)
	if err != nil {
		log.Err(err)
	}

	return leftDelegators
}

func (d delegationTask) getReturnedDelegators() []string {
	var returnedDelegators []string

	returnedDelegatorsQuery := `
		SELECT address 
		FROM delegation_history
		WHERE EXISTS (
			SELECT *
			FROM address_status
			WHERE 
				address_status.address = delegation_history.address
				AND
				address_status.status = 'leave'
		)
	`

	err := d.db.Select(&returnedDelegators, returnedDelegatorsQuery)
	if err != nil {
		log.Err(err)
	}

	return returnedDelegators
}

func RunDelegationTask(db *sqlx.DB) {
	task := newDelegationTask(db)

	task.getManagedChains()
	task.getNewDelegators()
	task.getLeftDelegators()
	task.getReturnedDelegators()

}
