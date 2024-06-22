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
	CREATE TABLE IF NOT EXISTS drive(
		id INTEGER PRIMARY KEY,
		idDriver INTEGER NOT NULL REFERENCES driver(id),
		sFrom TEXT NOT NULL,
		sTo TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_id ON drive(id);

	CREATE TABLE IF NOT EXISTS driver(
		id INTEGER PRIMARY KEY,
		sName TEXT NOT NULL,
		sPhone TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_sPhone ON driver(sPhone);

	`)

	if err != nil {
		return nil, fmt.Errorf("%s: error while executing : %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateDrive(idDriver int64, sFrom string, sTo string) error {
	const op = "storage.postgresql.CreateDrive"

	// ToDo: узнать, ппочему id не добавляется самостоятельно и сделать так, чтобы id подставлялся автоматически
	_, err := s.db.Exec("INSERT INTO drive(id, idDriver, sFrom, sTo) VALUES (1, $1, $2, $3);", idDriver, sFrom, sTo)

	if err != nil {
		return fmt.Errorf("%s: error while creating : %w", op, err)
	}

	return nil
}

func (s *Storage) GetDriverByDrive(idDrive int64) (int64, error) {
	const op = "storage.postgres.GetDriverByDrive"

	var idDriver int64
	err := s.db.QueryRow("SELECT idDriver FROM drive WHERE id = $1", idDrive).Scan(&idDriver)
	if err != nil {
		if err == pgx.ErrNoRows {
			return -1, fmt.Errorf("%s : error while getting driver : %w", op, pgx.ErrNoRows)
		}
		return -1, fmt.Errorf("%s : %w", op, err)
	}

	return idDriver, nil
}
