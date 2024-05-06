package db_connection

import (
	"WB0/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func DBConnect(cfg config.DBConnection) (*sql.DB, error) {

	var connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Успешное подключение к базе данных PostgreSQL")

	err = CreateTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
