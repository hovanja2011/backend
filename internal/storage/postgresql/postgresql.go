package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbPool struct {
	pool *pgxpool.Pool
}

func New(storagePath string) (*DbPool, error) {
	const op = "storage.postgresql.New"

	dbpool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	stmt := `
	CREATE TABLE IF NOT EXISTS drive(
		id INTEGER PRIMARY KEY,
		idDriver INTEGER NOT NULL,
		sFrom TEXT NOT NULL,
		sTo TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_phone ON drive(id);
	`
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s : Unable to begin transaction in pool: %v\n", op, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s : Error while transaction running: %v\n", op, err)
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	_, err = tx.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s : Query failed: %v\n", op, err)
		os.Exit(1)
	}

	return &DbPool{pool: dbpool}, nil
}

func (dbpool *DbPool) SaveDrive(idDriver int64, sFrom string, sTo string) (int64, error) {
	const op = "storage.postgresql.SaveDrive"

	tx, err := dbpool.pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s : Unable to begin transaction in pool: %v\n", op, err)
		return 0, err
	}
	var id int64
	row := tx.QueryRow(context.Background(), "INSERT INTO drive (idDriver, sFrom, sTo) VALUES ($1, $2, $3);", idDriver, sFrom, sTo)
	if err := row.Scan(&id); err != nil {
		fmt.Fprintf(os.Stderr, "%s : Error while inserting drive: %v\n", op, err)
		return 0, err
	}

	return id, nil

}
