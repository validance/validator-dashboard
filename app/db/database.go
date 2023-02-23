package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"sync"
	"validator-dashboard/app/config"
)

var once sync.Once
var db *sqlx.DB

func GetDB() *sqlx.DB {
	if db == nil {
		once.Do(newDB)
	}
	return db
}

func newDB() {
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

	database, dbOpenErr := sqlx.Connect("postgres", uri)

	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)

	if dbOpenErr != nil {
		log.Err(dbOpenErr).Msg("failed to open db")
	}

	db = database
}
