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
