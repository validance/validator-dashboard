package db

import "time"

type DelegationHistory struct {
	Id        int       `db:"id"`
	Address   string    `db:"address"`
	Validator string    `db:"validator"`
	Chain     string    `db:"chain"`
	Amount    float64   `db:"amount"`
	CreateDt  time.Time `db:"create_dt"`
	Label     string    `db:"label"`
	Status    string    `db:"status"`
}

type DelegationChanged struct {
	Address         string    `db:"address"`
	Validator       string    `db:"validator"`
	Chain           string    `db:"chain"`
	TodayAmount     string    `db:"today_amount"`
	YesterdayAmount string    `db:"yesterday_amount"`
	TodayDt         time.Time `db:"today_dt"`
	YesterdayDt     time.Time `db:"yesterday_dt"`
	Difference      float64   `db:"difference"`
	Label           string    `db:"label"`
	Status          string    `db:"status"`
}

type AddressStatus struct {
	Id        int       `db:"id"`
	Address   string    `db:"address"`
	Chain     string    `db:"chain"`
	Label     string    `db:"label"`
	Status    string    `db:"status"`
	CreatedDt time.Time `db:"create_dt"`
	UpdateDt  time.Time `db:"update_dt"`
}

type TokenPrice struct {
	Id        int       `db:"id"`
	Chain     string    `db:"chain"`
	Ticker    string    `db:"ticker"`
	Price     float64   `db:"price"`
	CreatedDt time.Time `db:"create_dt"`
}

type IncomeHistory struct {
	Id              int       `db:"id"`
	Address         string    `db:"address"`
	Chain           string    `db:"chain"`
	Reward          float64   `db:"reward"`
	Commission      float64   `db:"commission"`
	CreateDt        time.Time `db:"create_dt"`
	Price           float64   `db:"price"`
	RewardValue     float64   `db:"reward_value"`
	CommissionValue float64   `db:"commission_value"`
}

type DelegationSummaryHistory struct {
	Id                                            int       `db:"id"`
	Chain                                         string    `db:"chain"`
	Price                                         float64   `db:"price"`
	YesterdayDelegationAmountTotal                float64   `db:"yesterday_delegation_amount_total"`
	YesterdayDelegationValueTotal                 float64   `db:"yesterday_delegation_value_total"`
	YesterdayDelegationAmountB2B                  float64   `db:"yesterday_delegation_amount_b2b"`
	YesterdayDelegationValueB2B                   float64   `db:"yesterday_delegation_value_b2b"`
	YesterdayDelegationAmountB2C                  float64   `db:"yesterday_delegation_amount_b2c"`
	YesterdayDelegationValueB2C                   float64   `db:"yesterday_delegation_value_b2c"`
	YesterdayDelegationAmountUnknown              float64   `db:"yesterday_delegation_amount_unknown"`
	YesterdayDelegationValueUnknown               float64   `db:"yesterday_delegation_value_unknown"`
	TodayExistingIncreasedDelegationAmountTotal   float64   `db:"today_existing_increased_delegation_amount_total"`
	TodayExistingIncreasedDelegationValueTotal    float64   `db:"today_existing_increased_delegation_value_total"`
	TodayExistingIncreasedDelegationAmountB2B     float64   `db:"today_existing_increased_delegation_amount_b2b"`
	TodayExistingIncreasedDelegationValueB2B      float64   `db:"today_existing_increased_delegation_value_b2b"`
	TodayExistingIncreasedDelegationAmountB2C     float64   `db:"today_existing_increased_delegation_amount_b2c"`
	TodayExistingIncreasedDelegationValueB2C      float64   `db:"today_existing_increased_delegation_value_b2c"`
	TodayExistingIncreasedDelegationAmountUnknown float64   `db:"today_existing_increased_delegation_amount_unknown"`
	TodayExistingIncreasedDelegationValueUnknown  float64   `db:"today_existing_increased_delegation_value_unknown"`
	TodayNewIncreasedDelegationAmountTotal        float64   `db:"today_new_increased_delegation_amount_total"`
	TodayNewIncreasedDelegationValueTotal         float64   `db:"today_new_increased_delegation_value_total"`
	TodayNewIncreasedDelegationAmountB2B          float64   `db:"today_new_increased_delegation_amount_b2b"`
	TodayNewIncreasedDelegationValueB2B           float64   `db:"today_new_increased_delegation_value_b2b"`
	TodayNewIncreasedDelegationAmountB2C          float64   `db:"today_new_increased_delegation_amount_b2c"`
	TodayNewIncreasedDelegationValueB2C           float64   `db:"today_new_increased_delegation_value_b2c"`
	TodayNewIncreasedDelegationAmountUnknown      float64   `db:"today_new_increased_delegation_amount_unknown"`
	TodayNewIncreasedDelegationValueUnknown       float64   `db:"today_new_increased_delegation_value_unknown"`
	TodayReturnIncreasedDelegationAmountTotal     float64   `db:"today_return_increased_delegation_amount_total"`
	TodayReturnIncreasedDelegationValueTotal      float64   `db:"today_return_increased_delegation_value_total"`
	TodayReturnIncreasedDelegationAmountB2B       float64   `db:"today_return_increased_delegation_amount_b2b"`
	TodayReturnIncreasedDelegationValueB2B        float64   `db:"today_return_increased_delegation_value_b2b"`
	TodayReturnIncreasedDelegationAmountB2C       float64   `db:"today_return_increased_delegation_amount_b2c"`
	TodayReturnIncreasedDelegationValueB2C        float64   `db:"today_return_increased_delegation_value_b2c"`
	TodayReturnIncreasedDelegationAmountUnknown   float64   `db:"today_return_increased_delegation_amount_unknown"`
	TodayReturnIncreasedDelegationValueUnknown    float64   `db:"today_return_increased_delegation_value_unknown"`
	TodayExistingDecreasedDelegationAmountTotal   float64   `db:"today_existing_decreased_delegation_amount_total"`
	TodayExistingDecreasedDelegationValueTotal    float64   `db:"today_existing_decreased_delegation_value_total"`
	TodayExistingDecreasedDelegationAmountB2B     float64   `db:"today_existing_decreased_delegation_amount_b2b"`
	TodayExistingDecreasedDelegationValueB2B      float64   `db:"today_existing_decreased_delegation_value_b2b"`
	TodayExistingDecreasedDelegationAmountB2C     float64   `db:"today_existing_decreased_delegation_amount_b2c"`
	TodayExistingDecreasedDelegationValueB2C      float64   `db:"today_existing_decreased_delegation_value_b2c"`
	TodayExistingDecreasedDelegationAmountUnknown float64   `db:"today_existing_decreased_delegation_amount_unknown"`
	TodayExistingDecreasedDelegationValueUnknown  float64   `db:"today_existing_decreased_delegation_value_unknown"`
	TodayLeftDecreasedDelegationAmountTotal       float64   `db:"today_left_decreased_delegation_amount_total"`
	TodayLeftDecreasedDelegationValueTotal        float64   `db:"today_left_decreased_delegation_value_total"`
	TodayLeftDecreasedDelegationAmountB2B         float64   `db:"today_left_decreased_delegation_amount_b2b"`
	TodayLeftDecreasedDelegationValueB2B          float64   `db:"today_left_decreased_delegation_value_b2b"`
	TodayLeftDecreasedDelegationAmountB2C         float64   `db:"today_left_decreased_delegation_amount_b2c"`
	TodayLeftDecreasedDelegationValueB2C          float64   `db:"today_left_decreased_delegation_value_b2c"`
	TodayLeftDecreasedDelegationAmountUnknown     float64   `db:"today_left_decreased_delegation_amount_unknown"`
	TodayLeftDecreasedDelegationValueUnknown      float64   `db:"today_left_decreased_delegation_value_unknown"`
	TodayDelegationAmountTotal                    float64   `db:"today_delegation_amount_total"`
	TodayDelegationAmountB2B                      float64   `db:"today_delegation_amount_b2b"`
	TodayDelegationAmountB2C                      float64   `db:"today_delegation_amount_b2c"`
	TodayDelegationAmountUnknown                  float64   `db:"today_delegation_amount_unknown"`
	CreateDt                                      time.Time `db:"create_dt"`
}
