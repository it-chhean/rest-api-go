package store

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"email-api/models"
)

type EmailStore struct {
	db *sql.DB
}

func NewEmailStore(path string) (*EmailStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	store := &EmailStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return store, nil
}

func (s *EmailStore) Close() error {
	return s.db.Close()
}

func (s *EmailStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS emails (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT    NOT NULL UNIQUE
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *EmailStore) GetAll() ([]models.Email, error) {
	rows, err := s.db.Query("SELECT id, address FROM emails")
	if err != nil {
		return nil, fmt.Errorf("query emails: %w", err)
	}
	defer rows.Close()

	var emails []models.Email
	for rows.Next() {
		var e models.Email
		if err := rows.Scan(&e.ID, &e.Address); err != nil {
			return nil, fmt.Errorf("scan email row: %w", err)
		}
		emails = append(emails, e)
	}

	return emails, rows.Err()
}

func (s *EmailStore) GetByID(id int) (*models.Email, error) {
	var e models.Email
	err := s.db.QueryRow("SELECT id, address FROM emails WHERE id = ?", id).
		Scan(&e.ID, &e.Address)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query email by id: %w", err)
	}

	return &e, nil
}

func (s *EmailStore) Create(e *models.Email) error {
	result, err := s.db.Exec(
		"INSERT INTO emails (address) VALUES (?)", e.Address,
	)
	if err != nil {
		return fmt.Errorf("insert email: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	e.ID = int(id)
	return nil
}

func (s *EmailStore) Update(e *models.Email) error {
	result, err := s.db.Exec(
		"UPDATE emails SET address = ? WHERE id = ?", e.Address, e.ID,
	)
	if err != nil {
		return fmt.Errorf("update email: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (s *EmailStore) Delete(id int) error {
	result, err := s.db.Exec("DELETE FROM emails WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete email: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return models.ErrNotFound
	}

	return nil
}
