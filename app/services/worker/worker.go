package worker

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"sync"
	database "validator-dashboard/app/db"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services/client"
)

type worker struct {
	clients []client.Client
	db      *sql.DB
}

func spawnWorker(clients []client.Client, db *sql.DB) *worker {
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
		go w.spawnValidatorDelegationTask(&wg, c.ValidatorDelegations)
		go w.spawnValidatorIncomeTask(&wg, c.ValidatorIncome)
		go w.spawnGrantIncomeTask(&wg, c.GrantRewards)
	}

	wg.Wait()
}

func (w worker) spawnValidatorDelegationTask(wg *sync.WaitGroup, task func() (map[string]*models.Delegation, error)) {
	defer wg.Done()
	res, err := task()
	if err != nil {
		log.Err(err)
		return
	}

	query := `
		INSERT INTO delegation_history(address, validator, chain, amount) 
		VALUES ($1, $2, $3, $4)
	`

	for addr, delegation := range res {
		_, err := w.db.Exec(query, addr, delegation.Validator, delegation.Chain, delegation.Amount.String())
		if err != nil {
			log.Err(err)
		}
	}
}

func (w worker) spawnValidatorIncomeTask(wg *sync.WaitGroup, task func() (*models.ValidatorIncome, error)) {
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
}

func (w worker) spawnGrantIncomeTask(wg *sync.WaitGroup, task func() (map[string]*models.GrantReward, error)) {
	defer wg.Done()

	res, err := task()
	if err != nil {
		log.Err(err)
	}

	query := `
		INSERT INTO grant_reward_history(grant_address, validator, chain, reward)
		VALUES ($1, $2, $3, $4)
	`

	for grantAddr, reward := range res {
		_, exeErr := w.db.Exec(query, grantAddr, reward.Validator, reward.Chain, reward.Reward.String())
		if exeErr != nil {
			log.Err(exeErr)
		}
	}
}

func RunDbTask() error {
	log.Info().Msg("DB task running")

	clients, err := client.Initialize()
	if err != nil {
		return err
	}

	db, dbErr := database.New()
	if dbErr != nil {
		log.Err(dbErr)
	}

	defer db.Close()

	w := spawnWorker(clients, db)
	w.schedule()

	log.Info().Msg("DB task end")
	return nil
}
