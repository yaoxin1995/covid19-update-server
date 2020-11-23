package main

import (
	"covid19-update-service/model"
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
	updateServer, err := server.SetupServer(host, port)
	if err != nil {
		log.Fatalf("Could not start web server: %v", err)
	}

	go updateServer.Start()

	pollInterval, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL_MINUTES"))
	_ = rki.NewCovid19Poller(time.Duration(pollInterval) * time.Minute)

	done := make(chan os.Signal, 1)

	<-done
}
