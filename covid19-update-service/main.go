package main

import (
	"covid19-update-service/model"
	"covid19-update-service/server"
	"log"
	"os"
)

func init() {
	err := model.SetupDB("sqlite3", "storage.db")
	if err != nil {
		log.Fatalf("Could not setup database: %v", err)
	}
}

func main() {
	updateServer, err := server.SetupServer("localhost", "9005")
	if err != nil {
		log.Fatalf("Could not start web server: %v", err)
	}

	go updateServer.Start()

	done := make(chan os.Signal, 1)

	<-done
}
