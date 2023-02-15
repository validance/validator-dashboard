package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"validator-dashboard/app/config"
)

func New() *sql.DB {
	c := config.GetConfig()
	uri := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.Password,
		c.Database.DbName,
	)

	db, err := sql.Open("postgres", uri)

	if err != nil {
		fmt.Println(err)
	}

	connErr := db.Ping()
	if connErr == nil {
		fmt.Println("Connected!")
	}

	return db
}
