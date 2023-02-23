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

func (s *SummaryWorker) getCoinPrice(chain string) float64 {
	var price float64

	query := `
		SELECT DISTINCT price 
		FROM token_price
		WHERE
		    chain = $1 
		    AND
		    DATE(create_dt) = CURRENT_DATE + INTERVAL '-1 days'
	`

	err := s.db.Get(&price, query, chain)
	if err != nil {
		log.Err(err).Msg("cannot get token price")
		price = 0
	}

	return price
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
		Sum:     0,
	}
}

func (s *SummaryWorker) initSummaryWorker() {
	s.summary = make(map[string]*models.DelegationSummary)
	for _, c := range s.getManagedChains() {
		s.summary[c] = &models.DelegationSummary{
			YesterdayDelegationAmount:              newDelegationSummaryLabel(),
			TodayExistingIncreasedDelegationAmount: newDelegationSummaryLabel(),
			TodayNewIncreasedDelegationAmount:      newDelegationSummaryLabel(),
			TodayReturnIncreasedDelegationAmount:   newDelegationSummaryLabel(),
			TodayExistingDecreasedDelegationAmount: newDelegationSummaryLabel(),
			TodayLeftDecreasedDelegationAmount:     newDelegationSummaryLabel(),
		}
	}
}

func (s *SummaryWorker) sumUpDelegationValues(i *models.DelegationSummaryLabel) {
	i.Sum = i.B2B + i.B2C + i.Unknown
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
			s.sumUpDelegationValues(s.summary[d.Chain].TodayNewIncreasedDelegationAmount)

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
				s.sumUpDelegationValues(s.summary[d.Chain].TodayExistingIncreasedDelegationAmount)

			} else if d.Difference < 0 {
				switch d.Label {
				case "b2b":
					s.summary[d.Chain].TodayExistingDecreasedDelegationAmount.B2B += d.Difference
				case "b2c":
					s.summary[d.Chain].TodayExistingDecreasedDelegationAmount.B2C += d.Difference
				case "unknown":
					s.summary[d.Chain].TodayExistingDecreasedDelegationAmount.Unknown += d.Difference
				}
				s.sumUpDelegationValues(s.summary[d.Chain].TodayExistingDecreasedDelegationAmount)
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
			s.sumUpDelegationValues(s.summary[d.Chain].TodayLeftDecreasedDelegationAmount)

		case "return":
			switch d.Label {
			case "b2b":
				s.summary[d.Chain].TodayReturnIncreasedDelegationAmount.B2B += d.Difference
			case "b2c":
				s.summary[d.Chain].TodayReturnIncreasedDelegationAmount.B2C += d.Difference
			case "unknown":
				s.summary[d.Chain].TodayReturnIncreasedDelegationAmount.Unknown += d.Difference
			}
			s.sumUpDelegationValues(s.summary[d.Chain].TodayReturnIncreasedDelegationAmount)
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

func (s *SummaryWorker) runCreateDelegationSummaryTask() {
	createQuery := `
		INSERT INTO delegation_summary (
			chain, 
		                                
			yesterday_delegation_amount_total, 
			yesterday_delegation_amount_b2b, 
			yesterday_delegation_amount_b2c, 
			yesterday_delegation_amount_unknown,
		                                
			today_existing_increased_delegation_amount_total,
			today_existing_increased_delegation_amount_b2b,
			today_existing_increased_delegation_amount_b2c,
			today_existing_increased_delegation_amount_unknown,
		                                
			today_new_increased_delegation_amount_total,
			today_new_increased_delegation_amount_b2b,
			today_new_increased_delegation_amount_b2c,
			today_new_increased_delegation_amount_unknown,
		                                
			today_return_increased_delegation_amount_total,
			today_return_increased_delegation_amount_b2b,
			today_return_increased_delegation_amount_b2c,
			today_return_increased_delegation_amount_unknown,
		                                
			today_existing_decreased_delegation_amount_total,
			today_existing_decreased_delegation_amount_b2b,
			today_existing_decreased_delegation_amount_b2c,
			today_existing_decreased_delegation_amount_unknown,
		                                
			today_left_decreased_delegation_amount_total,
			today_left_decreased_delegation_amount_b2b,
			today_left_decreased_delegation_amount_b2c,
			today_left_decreased_delegation_amount_unknown
		)
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`
	for chain, summary := range s.summary {
		_, err := s.db.Exec(
			createQuery,

			chain,

			summary.YesterdayDelegationAmount.Sum,
			summary.YesterdayDelegationAmount.B2B,
			summary.YesterdayDelegationAmount.B2C,
			summary.YesterdayDelegationAmount.Unknown,

			summary.TodayExistingIncreasedDelegationAmount.Sum,
			summary.TodayExistingIncreasedDelegationAmount.B2B,
			summary.TodayExistingIncreasedDelegationAmount.B2C,
			summary.TodayExistingIncreasedDelegationAmount.Unknown,

			summary.TodayNewIncreasedDelegationAmount.Sum,
			summary.TodayNewIncreasedDelegationAmount.B2B,
			summary.TodayNewIncreasedDelegationAmount.B2C,
			summary.TodayNewIncreasedDelegationAmount.Unknown,

			summary.TodayReturnIncreasedDelegationAmount.Sum,
			summary.TodayReturnIncreasedDelegationAmount.B2B,
			summary.TodayReturnIncreasedDelegationAmount.B2C,
			summary.TodayReturnIncreasedDelegationAmount.Unknown,

			summary.TodayExistingDecreasedDelegationAmount.Sum,
			summary.TodayExistingDecreasedDelegationAmount.B2B,
			summary.TodayExistingDecreasedDelegationAmount.B2C,
			summary.TodayExistingDecreasedDelegationAmount.Unknown,

			summary.TodayLeftDecreasedDelegationAmount.Sum,
			summary.TodayLeftDecreasedDelegationAmount.B2B,
			summary.TodayLeftDecreasedDelegationAmount.B2C,
			summary.TodayLeftDecreasedDelegationAmount.Unknown,
		)
		log.Err(err).Msg("cannot create delegation summary")
	}
}

func (s *SummaryWorker) RunSummaryWorker() {
	s.initSummaryWorker()
	s.setPreviousDayDelegations()
	s.runDelegationChangedTask()
	s.runNewDelegatorTask()
	s.runLeftDelegatorTask()
	s.runReturnedDelegatorTask()
	s.runCreateDelegationSummaryTask()

	services.AddToCache("delegation_summary", s.summary)
}
