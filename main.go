package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s4lat/gokan/database"
	"github.com/s4lat/gokan/handlers"
	"github.com/s4lat/gokan/log"
)

func main() {
	// [INITIALIZING LOGGER]
	var logFile io.Writer
	if os.Getenv("GO_ENV") == "development" {
		logFile = os.Stdout
	} else {
		file, err := os.OpenFile("./shared/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logFile = file
	}
	logger := log.NewLogger(logFile)

	// [INITIALIZING DATABASE]
	DBURL := os.Getenv("DB_URL")
	if len(DBURL) == 0 {
		logger.Fatal("Environment variable 'DB_URL' is not set")
	}

	dbPool, err := pgxpool.New(context.Background(), DBURL)
	if err != nil {
		logger.Fatal(err)
	}
	// Checking that connection to db is established
	err = dbPool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	db := database.NewDB(dbPool)

	if os.Getenv("RECREATE_DB") == "1" {
		if err := db.System.RecreateAllTables(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}

	// [INITIALIZING HANDLERS AND SERVER]
	h := handlers.Handlers{DB: db, Log: logger}
	r := mux.NewRouter()
	r.HandleFunc("/", h.IndexHandler)

	s := http.Server{
		Addr:         os.Getenv("GOKAN_ADDR"),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      r,
	}

	logger.Info("Serving on " + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
