package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/processing-service/config"
	"github.com/processing-service/db"
	"github.com/processing-service/handler"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const appName = "processing-service"

func main() {
	config, err := config.ParseConfig(appName)
	if err != nil {
		log.Fatalf("parsing config; %v", err)
	}

	connectStr := fmt.Sprintf(config.DB.Tmpl, config.DB.Host, config.DB.Port, config.DB.Name, config.DB.User, appName)

	dbConn, err := db.CreateDB(connectStr, config.DB.ConnLifetime, config.DB.MaxIdleConns, config.DB.PoolSize)
	if err != nil {
		log.Fatalf("opening db connection: %v", err)
	}

	h, err := handler.CreateHTTPHandler(dbConn)
	if err != nil {
		log.Fatalf("creating http handler: %v", err)
	}

	listenErr := make(chan error, 1)
	server := &http.Server{
		Addr:         config.API.Port,
		ReadTimeout:  config.API.ReadTimeout,
		WriteTimeout: config.API.WriteTimeout,
		Handler:      h,
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
				lastTransactions, err := db.GetLastOddRecords(dbConn)
				if err != nil {
					log.Errorf("GetLastOddRecords: %v", err)
					continue
				}

				if len(lastTransactions) == 0 {
					continue
				}

				for _, t := range lastTransactions {
					log.Println("TRANSACTION", t)

					if err := db.PostProccessing(t, dbConn); err != nil {
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
