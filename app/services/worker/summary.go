package worker

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	database "validator-dashboard/app/db"
	"validator-dashboard/app/models"
	"validator-dashboard/app/services"
)

type SummaryWorker struct {
	db                 *sqlx.DB
	delegationChanged  []database.DelegationChanged
	newDelegators      []database.DelegationHistory
	leftDelegators     []database.DelegationHistory
	returnedDelegators []database.DelegationHistory
	summary            map[string]*models.DelegationSummary
}

func NewSummaryWorker(dt *DelegationTask) *SummaryWorker {
	return &SummaryWorker{
		dt.db,
		dt.delegationChanged,
		dt.newDelegators,
		dt.leftDelegators,
		dt.returnedDelegators,
		nil,
	}
}

func (s *SummaryWorker) getManagedChains() []string {
	var chains []string

	queryErr := s.db.Select(
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

func (s *SummaryWorker) setPreviousDayDelegations() {
	var pd []database.DelegationHistory

	query := `
		SELECT d.*, a.label, a.status
		FROM 
			address_status a JOIN delegation_history d
			ON a.address = d.address
		WHERE DATE(d.create_dt) = CURRENT_DATE + INTERVAL '-1 days'
	`

	err := s.db.Select(&pd, query)
	if err != nil {
		log.Err(err).Msg("failed to get previous day delegations")
	}

	for _, p := range pd {
		switch p.Label {
		case "b2b":
			s.summary[p.Chain].YesterdayDelegationAmount.B2B += p.Amount
		case "b2c":
			s.summary[p.Chain].YesterdayDelegationAmount.B2C += p.Amount
		case "unknown":
			s.summary[p.Chain].YesterdayDelegationAmount.Unknown += p.Amount
		}
	}
}

func newDelegationSummaryLabel() *models.DelegationSummaryLabel {
	return &models.DelegationSummaryLabel{
		B2B:     0,
		B2C:     0,
		Unknown: 0,
	}
}

func (s *SummaryWorker) initSummaryWorker() {
	s.summary = make(map[string]*models.DelegationSummary)
	for _, c := range s.getManagedChains() {
		s.summary[c] = &models.DelegationSummary{
			YesterdayDelegationAmount:              newDelegationSummaryLabel(),
			YesterdayDelegationValue:               newDelegationSummaryLabel(),
			TodayExistingIncreasedDelegationAmount: newDelegationSummaryLabel(),
			TodayExistingIncreasedDelegationValue:  newDelegationSummaryLabel(),
			TodayNewIncreasedDelegationAmount:      newDelegationSummaryLabel(),
			TodayNewIncreasedDelegationValue:       newDelegationSummaryLabel(),
			TodayReturnIncreasedDelegationAmount:   newDelegationSummaryLabel(),
			TodayReturnIncreasedDelegationValue:    newDelegationSummaryLabel(),
			TodayExistingDecreasedDelegationAmount: newDelegationSummaryLabel(),
			TodayExistingDecreasedDelegationValue:  newDelegationSummaryLabel(),
			TodayLeftDecreasedDelegationAmount:     newDelegationSummaryLabel(),
			TodayLeftDecreasedDelegationValue:      newDelegationSummaryLabel(),
		}
	}
}

func (s *SummaryWorker) runDelegationChangedTask() {
	for _, d := range s.delegationChanged {
		switch d.Status {
		case "new":
			switch d.Label {
			case "b2b":
				s.summary[d.Chain].TodayNewIncreasedDelegationAmount.B2B += d.Difference
			case "b2c":
				s.summary[d.Chain].TodayNewIncreasedDelegationAmount.B2C += d.Difference
			case "unknown":
				s.summary[d.Chain].TodayNewIncreasedDelegationAmount.Unknown += d.Difference
			}

		case "existing":
			if d.Difference > 0 {
				switch d.Label {
				case "b2b":
					s.summary[d.Chain].TodayExistingIncreasedDelegationAmount.B2B += d.Difference
				case "b2c":
					s.summary[d.Chain].TodayExistingIncreasedDelegationAmount.B2C += d.Difference
				case "unknown":
					s.summary[d.Chain].TodayExistingIncreasedDelegationAmount.Unknown += d.Difference
				}

			} else if d.Difference < 0 {
				switch d.Label {
				case "b2b":
					s.summary[d.Chain].TodayExistingDecreasedDelegationAmount.B2B += d.Difference
				case "b2c":
					s.summary[d.Chain].TodayExistingDecreasedDelegationAmount.B2C += d.Difference
				case "unknown":
					s.summary[d.Chain].TodayExistingDecreasedDelegationAmount.Unknown += d.Difference
				}
			}

		case "leave":
			switch d.Label {
			case "b2b":
				s.summary[d.Chain].TodayLeftDecreasedDelegationAmount.B2B += d.Difference
			case "b2c":
				s.summary[d.Chain].TodayLeftDecreasedDelegationAmount.B2C += d.Difference
			case "unknown":
				s.summary[d.Chain].TodayLeftDecreasedDelegationAmount.Unknown += d.Difference
			}

		case "return":
			switch d.Label {
			case "b2b":
				s.summary[d.Chain].TodayReturnIncreasedDelegationAmount.B2B += d.Difference
			case "b2c":
				s.summary[d.Chain].TodayReturnIncreasedDelegationAmount.B2C += d.Difference
			case "unknown":
				s.summary[d.Chain].TodayReturnIncreasedDelegationAmount.Unknown += d.Difference
			}
		}
	}
}

func (s *SummaryWorker) runNewDelegatorTask() {
	for _, n := range s.newDelegators {
		s.summary[n.Chain].TodayNewIncreasedDelegationAmount.Unknown += n.Amount
	}
}

func (s *SummaryWorker) runLeftDelegatorTask() {
	for _, l := range s.leftDelegators {
		switch l.Label {
		case "b2b":
			s.summary[l.Chain].TodayLeftDecreasedDelegationAmount.B2B += l.Amount
		case "b2c":
			s.summary[l.Chain].TodayLeftDecreasedDelegationAmount.B2C += l.Amount
		case "unknown":
			s.summary[l.Chain].TodayLeftDecreasedDelegationAmount.Unknown += l.Amount
		}
	}
}

func (s *SummaryWorker) runReturnedDelegatorTask() {
	for _, r := range s.returnedDelegators {
		switch r.Label {
		case "b2b":
			s.summary[r.Chain].TodayReturnIncreasedDelegationAmount.B2B += r.Amount
		case "b2c":
			s.summary[r.Chain].TodayReturnIncreasedDelegationAmount.B2C += r.Amount
		case "unknown":
			s.summary[r.Chain].TodayReturnIncreasedDelegationAmount.Unknown += r.Amount
		}
	}
}

func (s *SummaryWorker) RunSummaryWorker() {
	s.initSummaryWorker()
	s.setPreviousDayDelegations()
	s.runDelegationChangedTask()
	s.runNewDelegatorTask()
	s.runLeftDelegatorTask()
	s.runReturnedDelegatorTask()

	services.AddToCache("summary", s.summary)
}
