package postgres

import (
	"log"
	_"github.com/lib/pq"
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "Dilshod@2005"
	dbname   = "n9"
	port     = 5432
)

func Connection() *sql.DB{
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		log.Fatal(err)
	}
	return db
}