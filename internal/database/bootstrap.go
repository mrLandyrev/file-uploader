package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Bootstrap() (*sql.DB, error) {
	var db *sql.DB
	var err error

	if db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"192.168.0.105", 5432, "postgres", "aezakmi007", "postgres")); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS email (
		id UUID not null,
		filename text not null,
		send_date timestamp not null,
		upload_date timestamp not null,
		"to" text not null,
		cc text not null,
		subject text not null,
		content text not null,
		"from" text not null,
		primary key(id)
	)`); err != nil {
		return nil, err
	}

	return db, err
}
