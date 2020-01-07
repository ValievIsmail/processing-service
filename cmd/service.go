package main

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	cancelCtx = time.Second * 5
)

func dbGetLastOddRecords(db *sql.DB) (ts []Transaction, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cancelCtx)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select id, amount, state, src_type from api.get_last_odd_transactions()")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := Transaction{}

		if err := rows.Scan(&t.ID, &t.Amount, &t.State, &t.SrcType); err != nil {
			log.Errorf("dbGetLastOddRecords scan transaction %d: %v", t.ID, err)
			continue
		}

		ts = append(ts, t)
	}

	return ts, nil
}

// dbStateProccessing func
func dbStateProccessing(t Transaction, srcType string, db *sql.DB) error {
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

// dbPostProccessing func
func dbPostProccessing(t Transaction, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), cancelCtx)
	defer cancel()

	if _, err := db.ExecContext(ctx, "select api.post_processing($1, $2, $3, $4)", t.ID, t.Amount, t.State, t.SrcType); err != nil {
		return err
	}

	return nil
}
