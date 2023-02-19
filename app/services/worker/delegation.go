package worker

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	database "validator-dashboard/app/db"
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

func (d delegationTask) getNewDelegators() []database.DelegationHistory {
	var newDelegators []database.DelegationHistory

	newDelegatorQuery := `
		SELECT DISTINCT ON (d.address) d.*
		FROM 
			delegation_history d LEFT JOIN address_status a
			ON d.address = a.address
		WHERE a IS NULL
	`

	err := d.db.Select(&newDelegators, newDelegatorQuery)
	if err != nil {
		log.Err(err)
	}

	return newDelegators
}

func (d delegationTask) getLeftDelegators() []database.DelegationHistory {
	var leftDelegators []database.DelegationHistory

	leftDelegatorsQuery := `
		SELECT yesterday.*
		FROM 
			(
				SELECT *
				FROM delegation_history
				WHERE DATE(delegation_history.create_dt) = CURRENT_DATE + INTERVAL '-1 DAYS'
			) yesterday LEFT JOIN
			(
				SELECT *
				FROM delegation_history
				WHERE DATE(delegation_history.create_dt) = CURRENT_DATE
			) today
			ON yesterday.address = today.address
		WHERE today is NULL
	`

	err := d.db.Select(&leftDelegators, leftDelegatorsQuery)
	if err != nil {
		log.Err(err)
	}

	return leftDelegators
}

func (d delegationTask) getReturnedDelegators() []database.DelegationHistory {
	var returnedDelegators []database.DelegationHistory

	returnedDelegatorsQuery := `
		SELECT DISTINCT ON (d.address) d.*
		FROM 
			delegation_history d LEFT JOIN address_status a
			ON d.address = a.address
		WHERE a.status = 'leave'
	`

	err := d.db.Select(&returnedDelegators, returnedDelegatorsQuery)
	if err != nil {
		log.Err(err)
	}

	return returnedDelegators
}

func (d delegationTask) getDelegationChanged() []database.DelegationChanged {
	var delegationChanged []database.DelegationChanged

	query := `
		SELECT 	
		    today.address, 
			today.validator, 
			today.chain, 
			today.amount as today_amount, 
			yesterday.amount as yesterday_amount, 
			ROUND(cast(today.amount as float)::numeric - cast(yesterday.amount as float)::numeric, 5) as difference
		FROM (
				SELECT *
				FROM delegation_history
				WHERE DATE(delegation_history.create_dt) = CURRENT_DATE
			) today JOIN
			(
				SELECT *
				FROM delegation_history
				WHERE DATE(delegation_history.create_dt) = CURRENT_DATE + INTERVAL '-1 DAYS'
			) yesterday
		ON 
			today.address = yesterday.address
			AND	
			today.validator = yesterday.validator 
			AND
			today.amount != yesterday.amount
	`

	err := d.db.Select(&delegationChanged, query)
	if err != nil {
		log.Err(err)
	}

	return delegationChanged
}

func RunDelegationTask(db *sqlx.DB) {
	task := newDelegationTask(db)

	//fmt.Println(task.getManagedChains())
	for _, val := range task.getNewDelegators() {
		fmt.Println(val)
	}
	for _, val := range task.getLeftDelegators() {
		fmt.Println(val)
	}
	fmt.Println(task.getReturnedDelegators())

	fmt.Println("delegation changed: ")
	for _, val := range task.getDelegationChanged() {
		fmt.Println(val)
	}
}
