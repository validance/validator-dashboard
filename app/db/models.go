package db

import "time"

type DelegationHistory struct {
	Id        int       `db:"id"`
	Address   string    `db:"address"`
	Validator string    `db:"validator"`
	Chain     string    `db:"chain"`
	Amount    string    `db:"amount"`
	CreatedDt time.Time `db:"create_dt"`
}

type DelegationChanged struct {
	Address         string  `db:"address"`
	Validator       string  `db:"validator"`
	Chain           string  `db:"chain"`
	TodayAmount     string  `db:"today_amount"`
	YesterdayAmount string  `db:"yesterday_amount"`
	Difference      float64 `db:"difference"`
}
