package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"validator-dashboard/app/config"
)

func New() (*sql.DB, error) {
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

	db, dbOpenErr := sql.Open("postgres", uri)

	if dbOpenErr != nil {
		return nil, dbOpenErr
	}

	dbConnErr := db.Ping()

	if dbConnErr != nil {
		return nil, dbConnErr
	}

	return db, nil
}
