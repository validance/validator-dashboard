package worker

import (
	"sync"
	database "validator-dashboard/app/db"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services/client"

	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type worker struct {
	clients []client.Client
	db      *sqlx.DB
}

func spawnHistoryWorker(clients []client.Client, db *sqlx.DB) *worker {
	return &worker{
		clients,
		db,
	}
}

func (w worker) RunHistoryWorker() {
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
		log.Err(err).Msg("failed to spawn validator delegation history task")
		return
	}

	chain := ""

	query := `
		INSERT INTO delegation_history(address, validator, chain, amount) 
		VALUES ($1, $2, $3, $4)
	`

	for addr, delegation := range res {
		_, err := w.db.Exec(query, addr, delegation.Validator, delegation.Chain, delegation.Amount)
		if err != nil {
			log.Err(err).Msg("failed to insert delegation history")
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
		log.Err(err).Msg("failed to spawn validator income history task")
		return
	}

	query := `
		INSERT INTO income_history(address, chain, reward, commission)
		VALUES ($1, $2, $3, $4)
	`

	_, exeErr := w.db.Exec(query, res.Validator, res.Chain, res.Reward, res.Commission)

	if exeErr != nil {
		log.Err(exeErr).Msg("failed to create validtor income history")
	}

	if res.Chain != "" {
		log.Info().Msgf("validator income task finished: %s", res.Chain)
	}
}

func (w worker) spawnGrantIncomeHistoryTask(wg *sync.WaitGroup, task func() (map[string]*models.GrantReward, error)) {
	defer wg.Done()

	res, err := task()
	if err != nil {
		log.Err(err).Msg("failed to spawn grant income history task")
	}

	chain := ""

	query := `
		INSERT INTO grant_reward_history(grant_address, validator, chain, reward)
		VALUES ($1, $2, $3, $4)
	`

	for grantAddr, reward := range res {
		_, exeErr := w.db.Exec(query, grantAddr, reward.Validator, reward.Chain, reward.Reward)
		chain = reward.Chain
		if exeErr != nil {
			log.Err(exeErr).Msg("failed to create grant income history")
		}
	}

	if chain != "" {
		log.Info().Msgf("validator reward task finished: %s", chain)
	}
}

// query on chain data and insert those in to db
func run() error {
	clients := client.Initialize()

	db := database.GetDB()

	// pipelining tasks
	log.Info().Msg("History task running")
	hw := spawnHistoryWorker(clients, db)
	hw.RunHistoryWorker()
	log.Info().Msg("History task end")

	log.Info().Msg("Delegation status task running")
	dt := NewDelegationStatusTask(db)
	dt.RunDelegationStatusTask()
	log.Info().Msg("Delegation status task end")

	log.Info().Msg("TokenPrice task running")
	tp := NewTokenPriceTask(db)
	tp.RunTokenPriceTask()
	log.Info().Msg("TokenPrice task end")

	log.Info().Msg("Summary task running")
	sw := NewSummaryWorker(dt)
	sw.RunSummaryWorker()
	log.Info().Msg("Summary task end")

	return nil
}

func runJob() {
	<-gocron.Start()
}

func RegisterCron(t string) {
	gocron.Every(1).Day().At(t).Do(run)
	go runJob()
}
