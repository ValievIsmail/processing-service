package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/processing-service/models"

	log "github.com/sirupsen/logrus"
)

var (
	cancelCtx = time.Second * 5
)

// CreateDB func
func CreateDB(connectStr string, connLife time.Duration, maxIdle, poolSize int) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(poolSize)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(connLife)

	return db, nil
}

// GetLastOddRecords func
func GetLastOddRecords(db *sql.DB) (ts []models.Transaction, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cancelCtx)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select id, amount, state, src_type from api.get_last_odd_transactions()")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := models.Transaction{}

		if err := rows.Scan(&t.ID, &t.Amount, &t.State, &t.SrcType); err != nil {
			log.Errorf("dbGetLastOddRecords scan transaction %d: %v", t.ID, err)
			continue
		}

		ts = append(ts, t)
	}

	return ts, nil
}

// StateProccessing func
func StateProccessing(t models.Transaction, srcType string, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), cancelCtx)
	defer cancel()

	var query string

	if t.State == "win" {
		query = `select api.win_processing($1, $2, $3)`
	} else {
		query = `select api.lost_processing($1, $2, $3)`
	}

	if _, err := db.ExecContext(ctx, query, t.ID, t.Amount, srcType); err != nil {
		return err
	}

	return nil
}

// PostProccessing func
func PostProccessing(t models.Transaction, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), cancelCtx)
	defer cancel()

	if _, err := db.ExecContext(ctx, "select api.post_processing($1, $2, $3, $4)", t.ID, t.Amount, t.State, t.SrcType); err != nil {
		return err
	}

	return nil
}
