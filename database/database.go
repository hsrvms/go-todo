package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PgStorage struct {
	db *sql.DB
}

// NewPgStorage initializes a new PgStorage with the given connection string.
//
// Parameters:
//
//	connStr - the connection string for the PostgreSQL database.
//
// Returns:
//
//	*PgStorage - a pointer to the PgStorage instance.
func NewPgStorage(connStr string) *PgStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return &PgStorage{db: db}
}

// Init initializes the PgStorage by creating the users and tasks tables.
func (s *PgStorage) Init() (*sql.DB, error) {
	if err := s.createUsersTable(); err != nil {
		return nil, err
	}

	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *PgStorage) createUsersTable() error {
	if s == nil || s.db == nil {
		return errors.New("nil receiver or nil db connection")
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table %v", err)
	}

	return nil
}

func (s *PgStorage) createTasksTable() error {
	if s == nil || s.db == nil {
		return errors.New("nil receiver or nil db connection")
	}

	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL,
			status BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %v", err)
	}

	return nil
}
