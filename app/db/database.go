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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.DbName,
		5,
	)

	db, err := sql.Open("postgres", uri)

	if err != nil {
		fmt.Println(err)
	}

	connErr := db.Ping()

	if connErr != nil {
		fmt.Println(connErr)
		return nil
	} else {
		fmt.Println("Database Connected!")
	}

	return db
}
