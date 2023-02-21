package db

import "time"

type DelegationHistory struct {
	Id        int       `db:"id"`
	Address   string    `db:"address"`
	Validator string    `db:"validator"`
	Chain     string    `db:"chain"`
	Amount    string    `db:"amount"`
	CreateDt  time.Time `db:"create_dt"`
}

type DelegationChanged struct {
	Address         string  `db:"address"`
	Validator       string  `db:"validator"`
	Chain           string  `db:"chain"`
	TodayAmount     string  `db:"today_amount"`
	YesterdayAmount string  `db:"yesterday_amount"`
	Difference      float64 `db:"difference"`
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
