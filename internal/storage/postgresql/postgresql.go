package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: error while open connect : %w", op, err)
	}
	// defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS driver(
		id INTEGER PRIMARY KEY,
		sPhone TEXT NOT NULL UNIQUE,
		sName TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_sPhone ON driver(sPhone);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: error while executing : %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateDriver(driverPhone string, driverName string) error {
	const op = "storage.postgresql.createdriver"

	// ToDo: узнать, ппочему id не добавляется самостоятельно и сделать так, чтобы id подставлялся автоматически
	_, err := s.db.Exec("INSERT INTO driver(id, sPhone, sName) VALUES (1, $1, $2);", driverPhone, driverName)

	if err != nil {
		return fmt.Errorf("%s: error while creating : %w", op, err)
	}

	return nil
}

func (s *Storage) GetDriver(driverPhone string) (string, error) {
	const op = "storage.postgres.getdriver"

	var driverName string
	err := s.db.QueryRow("SELECT sName FROM driver WHERE sPhone = $1", driverPhone).Scan(&driverName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("%s : error while getting driver : %w", op, pgx.ErrNoRows)
		}
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return driverName, nil
}
