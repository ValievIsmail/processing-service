package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const appName = "processing-service"

func main() {
	config, err := parseConfig(appName)
	if err != nil {
		log.Fatalf("parsing config; %v", err)
	}

	connectStr := fmt.Sprintf(config.DB.Tmpl, config.DB.Host, config.DB.Port, config.DB.Name, config.DB.User, appName)

	db, err := createDB(connectStr, config.DB.ConnLifetime, config.DB.MaxIdleConns, config.DB.PoolSize)
	if err != nil {
		log.Errorf("opening connection: %v", err)
	}
	log.RegisterExitHandler(func() {
		db.Close()
	})

	handler, err := createHTTPHandler(db)
	if err != nil {
		log.Fatalf("creating http handler: %v", err)
	}

	listenErr := make(chan error, 1)
	server := &http.Server{
		Addr:         config.API.Port,
		ReadTimeout:  config.API.ReadTimeout,
		WriteTimeout: config.API.WriteTimeout,
		Handler:      handler,
	}

	go func() {
		log.Println("PROCESSING-SERVICE STARTED")
		listenErr <- server.ListenAndServe()
	}()

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(config.API.ProccesingTime):
				lastTransactions, err := dbGetLastOddRecords(db)
				if err != nil {
					log.Errorf("dbGetLastOddRecords: %v", err)
					continue
				}

				if len(lastTransactions) == 0 {
					continue
				}

				for _, t := range lastTransactions {
					log.Println("TRANSACTION", t)

					if err := dbPostProccessing(t, db); err != nil {
						log.Errorf("dbStateProccessing post processing with id %d: %v", t.ID, err)
						continue
					}
				}
			case <-quit:
				return
			}
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		log.Fatal(err)
	case <-osSignals:
		server.SetKeepAlivesEnabled(false)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("STOP APP")
		log.Exit(0)
	}
}

func createDB(connectStr string, connLife time.Duration, maxIdle, poolSize int) (db *sql.DB, err error) {
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
