package main

import (
	"covid19-update-service/model"
	"covid19-update-service/notifier"
	"covid19-update-service/rki"
	"covid19-update-service/server"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	dbType := os.Getenv("DB_TYPE")
	dbSource := os.Getenv("DB_SOURCE")
	err := model.SetupDB(dbType, dbSource)
	if err != nil {
		log.Fatalf("Could not setup database: %v", err)
	}
}

func main() {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	respAPI, err := server.SetupServer(host, port)
	if err != nil {
		log.Fatalf("Could not start web server: %v", err)
	}

	go respAPI.Start()

	pollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL_MINUTES"))
	if err != nil {
		log.Fatalf("could not load poll interval: %v", err)
	}

	c := make(chan model.Covid19Region, 500)
	_ = rki.NewRegionUpdater(time.Duration(pollInterval)*time.Minute, c)
	_ = notifier.NewCovid19Notifier(c)

	done := make(chan os.Signal, 1)
	<-done
}
