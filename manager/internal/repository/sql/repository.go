package sql

import (
	"database/sql"
	"errors"
	"sync"

	"manager/internal/domain"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName     = "sqlite3"
	dataSourceName = "./db/manager.db"
)

type Repository struct {
	db   *sql.DB
	lock sync.RWMutex
}

func NewRepository() *Repository {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic("failed to establish database connection: " + err.Error())
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS crack_hash (
	    request_id UUID NOT NULL PRIMARY KEY,
		hash VARCHAR NOT NULL,
		max_length INTEGER NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		_ = db.Close()
		panic("failed to create crack_hash table: " + err.Error())
	}

	return &Repository{
		db:   db,
		lock: sync.RWMutex{},
	}
}

func (r *Repository) Save(requestId string, hash string, maxLength int) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if requestId == "" || hash == "" {
		return errors.New("request_id and hash cannot be empty")
	}

	id, err := uuid.Parse(requestId)
	if err != nil {
		return errors.New("invalid request_id format: " + err.Error())
	}

	query := `INSERT OR REPLACE INTO crack_hash (request_id, hash, max_length) VALUES (?, ?, ?)`
	_, err = r.db.Exec(query, id, hash, maxLength)
	if err != nil {
		return errors.New("failed to save crack hash: " + err.Error())
	}

	return nil
}

func (r *Repository) Remove(requestId string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if requestId == "" {
		return errors.New("request_id cannot be empty")
	}

	id, err := uuid.Parse(requestId)
	if err != nil {
		return errors.New("invalid request_id format: " + err.Error())
	}

	query := `DELETE FROM crack_hash WHERE request_id = ?`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return errors.New("failed to remove crack hash: " + err.Error())
	}

	return nil
}

func (r *Repository) GetAll() ([]domain.CrackHash, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	query := `SELECT request_id, hash, max_length FROM crack_hash`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.New("failed to query crack hashes: " + err.Error())
	}
	defer func() {
		_ = rows.Close()
	}()

	var results []domain.CrackHash
	for rows.Next() {
		var ch domain.CrackHash
		if err = rows.Scan(&ch.RequestId, &ch.Hash, &ch.MaxLength); err != nil {
			return nil, errors.New("failed to scan crack hash row: " + err.Error())
		}
		results = append(results, ch)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("error occurred while iterating rows: " + err.Error())
	}

	return results, nil
}

func (r *Repository) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.db.Close(); err != nil {
		return errors.New("failed to close database connection: " + err.Error())
	}

	return nil
}
