package worker

import (
	"fmt"
	"sync"
	database "validator-dashboard/app/db"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services/client"
)

type worker struct {
	clients []client.Client
}

func spawnWorker(clients []client.Client) *worker {
	return &worker{clients}
}

func (w worker) schedule() {
	var wg sync.WaitGroup

	tasksNum := 1
	wg.Add(len(w.clients) * tasksNum)

	for _, c := range w.clients {
		//go w.spawnValidatorDelegationTask(&wg, c.ValidatorDelegations)
		go w.spawnValidatorIncomeTask(&wg, c.ValidatorIncome)
	}

	wg.Wait()
}

func (w worker) spawnValidatorDelegationTask(wg *sync.WaitGroup, task func() (map[string]*models.Delegation, error)) {
	defer wg.Done()
	res, err := task()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, dbErr := database.New()
	if dbErr != nil {
		fmt.Println(dbErr)
		return
	}

	defer db.Close()

	query := `
		INSERT INTO delegation_history(address, validator, chain, amount) 
		VALUES ($1, $2, $3, $4)
	`

	for addr, delegation := range res {
		_, err := db.Exec(query, addr, delegation.Validator, delegation.Chain, delegation.Amount.String())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (w worker) spawnValidatorIncomeTask(wg *sync.WaitGroup, task func() (*models.ValidatorIncome, error)) {
	defer wg.Done()
	res, err := task()

	if err != nil {
		fmt.Println(err)
		return
	}

	db, dbErr := database.New()
	if dbErr != nil {
		fmt.Println(dbErr)
		return
	}

	defer db.Close()

	query := `
		INSERT INTO income_history(address, chain, reward, commission)
		VALUES ($1, $2, $3, $4)
	`
	_, exeErr := db.Exec(query, res.Validator, res.Chain, res.Reward.String(), res.Commission.String())

	if exeErr != nil {
		fmt.Println(exeErr)
	}
}

func Run() error {
	clients, err := client.Initialize()
	if err != nil {
		return err
	}

	w := spawnWorker(clients)
	w.schedule()

	return nil
}
