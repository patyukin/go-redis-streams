package storage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return &Storage{db: db}, nil
}

func (s *Storage) Save() error {
	return nil
}

func (s *Storage) GetById(id int) error {
	return nil
}

func (s *Storage) GetAllIds() error {
	return nil
}
