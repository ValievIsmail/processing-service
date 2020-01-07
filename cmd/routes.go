package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func createHTTPHandler(db *sql.DB) (http.Handler, error) {
	mux := chi.NewMux()

	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(200)))
	})

	mux.Route("/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8;"))
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           30,
		}).Handler)

		r.Post("/process", func(w http.ResponseWriter, r *http.Request) {
			t := Transaction{}

			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, "process handler decode err", http.StatusBadRequest)
				return
			}

			if !validateTransaction(t) {
				http.Error(w, "process incorrect data err", http.StatusBadRequest)
				return
			}

			srcType := r.Header.Get("Source-Type")

			if err := dbStateProccessing(t, srcType, db); err != nil {
				http.Error(w, fmt.Sprintf("dbStateProccessing err with transaction %d: %v", t.ID, err), http.StatusInternalServerError)
				return
			}
		})
	})

	return mux, nil
}

func validateTransaction(t Transaction) bool {
	if t.ID == 0 || t.State != "win" && t.State != "lost" {
		return false
	}

	return true
}
