package worker

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"sync"
	database "validator-dashboard/app/db"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services/client"
)

type worker struct {
	clients []client.Client
	db      *sqlx.DB
}

func spawnWorker(clients []client.Client, db *sqlx.DB) *worker {
	return &worker{
		clients,
		db,
	}
}

func (w worker) schedule() {
	var wg sync.WaitGroup

	// set number of task to be joined by goroutine
	tasksNum := 3
	wg.Add(len(w.clients) * tasksNum)

	for _, c := range w.clients {
		go w.spawnValidatorDelegationHistoryTask(&wg, c.ValidatorDelegations)
		go w.spawnValidatorIncomeHistoryTask(&wg, c.ValidatorIncome)
		go w.spawnGrantIncomeHistoryTask(&wg, c.GrantRewards)
	}

	wg.Wait()
}

func (w worker) spawnValidatorDelegationHistoryTask(wg *sync.WaitGroup, task func() (map[string]*models.Delegation, error)) {
	defer wg.Done()
	res, err := task()
	if err != nil {
		log.Err(err)
		return
	}

	chain := ""

	query := `
		INSERT INTO delegation_history(address, validator, chain, amount) 
		VALUES ($1, $2, $3, $4)
	`

	addrStatusQuery := `
		INSERT INTO address_status (address, chain)
		VALUES ($1, $2)
	`

	for addr, delegation := range res {
		_, err := w.db.Exec(query, addr, delegation.Validator, delegation.Chain, delegation.Amount.String())
		_, e := w.db.Exec(addrStatusQuery, addr, delegation.Chain)
		_ = e
		if err != nil {
			log.Err(err)
		}
		chain = delegation.Chain
	}

	if chain != "" {
		log.Info().Msgf("validator delegation task finished: %s", chain)
	}
}

func (w worker) spawnValidatorIncomeHistoryTask(wg *sync.WaitGroup, task func() (*models.ValidatorIncome, error)) {
	defer wg.Done()
	res, err := task()

	if err != nil {
		log.Err(err)
		return
	}

	query := `
		INSERT INTO income_history(address, chain, reward, commission)
		VALUES ($1, $2, $3, $4)
	`

	_, exeErr := w.db.Exec(query, res.Validator, res.Chain, res.Reward.String(), res.Commission.String())

	if exeErr != nil {
		log.Err(exeErr)
	}

	if res.Chain != "" {
		log.Info().Msgf("validator income task finished: %s", res.Chain)
	}
}

func (w worker) spawnGrantIncomeHistoryTask(wg *sync.WaitGroup, task func() (map[string]*models.GrantReward, error)) {
	defer wg.Done()

	res, err := task()
	if err != nil {
		log.Err(err)
	}

	chain := ""

	query := `
		INSERT INTO grant_reward_history(grant_address, validator, chain, reward)
		VALUES ($1, $2, $3, $4)
	`

	for grantAddr, reward := range res {
		_, exeErr := w.db.Exec(query, grantAddr, reward.Validator, reward.Chain, reward.Reward.String())
		chain = reward.Chain
		if exeErr != nil {
			log.Err(exeErr)
		}
	}

	if chain != "" {
		log.Info().Msgf("validator reward task finished: %s", chain)
	}
}

// Run query on chain data and insert those in to db
func Run() error {
	log.Info().Msg("DB task running")

	clients, err := client.Initialize()
	if err != nil {
		return err
	}

	db, dbErr := database.New()
	if dbErr != nil {
		log.Err(dbErr)
		return dbErr
	}

	defer db.Close()

	w := spawnWorker(clients, db)
	w.schedule()
	//_ = w

	RunDelegationTask(db)

	log.Info().Msg("DB task end")
	return nil
}
