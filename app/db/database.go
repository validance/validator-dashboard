package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"validator-dashboard/app/config"
)

func New() (*sqlx.DB, error) {
	c := config.GetConfig()
	uri := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.DbName,
		5,
	)

	db, dbOpenErr := sqlx.Connect("postgres", uri)

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if dbOpenErr != nil {
		return nil, dbOpenErr
	}

	return db, nil
}
