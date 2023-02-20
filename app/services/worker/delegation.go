package worker

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"sync"
	database "validator-dashboard/app/db"
)

type DelegationTask struct {
	db                 *sqlx.DB
	delegationChanged  []database.DelegationChanged
	newDelegators      []database.DelegationHistory
	leftDelegators     []database.DelegationHistory
	returnedDelegators []database.DelegationHistory
}

func NewDelegationTask(db *sqlx.DB) *DelegationTask {

	return &DelegationTask{
		db,
		nil,
		nil,
		nil,
		nil,
	}
}

func (d *DelegationTask) getManagedChains() []string {
	var chains []string

	queryErr := d.db.Select(
		&chains,
		`
			SELECT DISTINCT chain
			FROM delegation_history
		`,
	)

	if queryErr != nil {
		log.Err(queryErr).Msg("failed to get managed chains")
	}

	return chains
}

func (d *DelegationTask) getNewDelegators() []database.DelegationHistory {
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
		log.Err(err).Msg("failed to get new delegators")
	}
	return newDelegators
}

func (d *DelegationTask) getLeftDelegators() []database.DelegationHistory {
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
		log.Err(err).Msg("failed to get left delegators")
	}

	return leftDelegators
}

func (d *DelegationTask) getReturnedDelegators() []database.DelegationHistory {
	var returnedDelegators []database.DelegationHistory

	returnedDelegatorsQuery := `
		SELECT DISTINCT ON (d.address) d.*
		FROM 
			delegation_history d LEFT JOIN address_status a
			ON d.address = a.address
		WHERE 
		    a.status = 'leave'
			AND
			DATE(d.create_dt) = CURRENT_DATE 
	`

	err := d.db.Select(&returnedDelegators, returnedDelegatorsQuery)
	if err != nil {
		log.Err(err).Msg("failed to get returned delegators")
	}

	return returnedDelegators
}

func (d *DelegationTask) getDelegationChanged() []database.DelegationChanged {
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
		log.Err(err).Msg("failed to get delegation changed")
	}

	return delegationChanged
}

func (d *DelegationTask) createNewDelegators(dhs []database.DelegationHistory) {
	createQuery := `
		INSERT INTO address_status(address, chain)
		VALUES ($1, $2)
	`

	for _, dh := range dhs {
		_, err := d.db.Exec(createQuery, dh.Address, dh.Chain)
		if err != nil {
			log.Err(err).Msg("failed to create new delegator status")
		}
	}
}

func (d *DelegationTask) updateLeftDelegators(dhs []database.DelegationHistory) {
	updateQuery := `
 		UPDATE address_status
		SET
		    status = 'leave',
			update_dt = CURRENT_TIMESTAMP
		WHERE address = $1
	`

	for _, dh := range dhs {
		_, err := d.db.Exec(updateQuery, dh.Address)
		if err != nil {
			log.Err(err).Msg("failed to update left delegators")
		}
	}
}

func (d *DelegationTask) updateReturnedDelegators(dhs []database.DelegationHistory) {
	updateQuery := `
		UPDATE address_status
		SET
		    status = 'return',
			update_dt = CURRENT_TIMESTAMP
		WHERE address = $1
	`

	for _, dh := range dhs {
		_, err := d.db.Exec(updateQuery, dh.Address)
		if err != nil {
			log.Err(err).Msg("error on updating returned delegators")
		}
	}
}

func (d *DelegationTask) updateExistingDelegators() {
	updateQuery := `
		UPDATE address_status
		SET
		    status = 'existing',
			update_dt = CURRENT_TIMESTAMP
		WHERE
		    status = 'new'
	`

	_, err := d.db.Exec(updateQuery)
	if err != nil {
		log.Err(err).Msg("error on updating existing delegators")
	}
}

func (d *DelegationTask) RunDelegationTask() {
	tasksNum := 4
	var wg sync.WaitGroup

	wg.Add(tasksNum)

	go func() {
		defer wg.Done()
		changedDelegators := d.getDelegationChanged()
		d.delegationChanged = changedDelegators
		d.updateExistingDelegators()
	}()

	go func() {
		defer wg.Done()
		newDelegators := d.getNewDelegators()
		d.newDelegators = newDelegators
		d.createNewDelegators(newDelegators)
	}()

	go func() {
		defer wg.Done()
		leftDelegators := d.getLeftDelegators()
		d.leftDelegators = leftDelegators
		d.updateLeftDelegators(leftDelegators)
	}()

	go func() {
		defer wg.Done()
		returnDelegators := d.getReturnedDelegators()
		d.returnedDelegators = returnDelegators
		d.updateReturnedDelegators(returnDelegators)
	}()

	wg.Wait()
}
