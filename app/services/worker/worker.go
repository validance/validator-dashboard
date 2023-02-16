package worker

import (
	"fmt"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services/client"
)

type worker struct {
	clients []client.Client
}

func spawnWorker(clients []client.Client) *worker {
	return &worker{clients}
}

func (w worker) schedule() error {
	var vdt []func() (map[string]*models.Delegation, error)
	var vit []func() (*models.ValidatorIncome, error)

	for _, c := range w.clients {
		vdt = append(vdt, c.ValidatorDelegations)
	}

	for _, c := range w.clients {
		vit = append(vit, c.ValidatorIncome)
	}

	vdtErr := w.spawnValidatorDelegationTask(vdt)

	if vdtErr != nil {
		return vdtErr
	}

	vitErr := w.spawnValidatorIncomeTask(vit)
	if vitErr != nil {
		return vitErr
	}

	return nil
}

func (w worker) spawnValidatorDelegationTask(fns []func() (map[string]*models.Delegation, error)) error {
	for _, f := range fns {
		f := f
		go func() {
			res, err := f()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	// TODO: put data into database
	return err
}

func (w worker) spawnValidatorIncomeTask(fns []func() (*models.ValidatorIncome, error)) error {
	result, err := f()
	// TODO: put data into database
	_ = result
	return err
}

func Run() error {
	clients, err := client.Initialize()
	if err != nil {
		return err
	}

	w := spawnWorker(clients)
	return w.schedule()
}
