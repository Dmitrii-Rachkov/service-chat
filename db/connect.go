package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDB - соединение с базой данных postgres
func ConnectDB(host, port, user, password, name string, conn int) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, name))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(conn)
	return db, nil
}
