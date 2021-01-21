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
	// DB Setup
	dbType, ok := os.LookupEnv("DB_TYPE")
	if !ok {
		log.Fatalf("DB_TYPE missing")
	}
	dbSource, ok := os.LookupEnv("DB_SOURCE")
	if !ok {
		log.Fatalf("DB_SOURCE missing")
	}
	err := model.SetupDB(dbType, dbSource)
	if err != nil {
		log.Fatalf("Could not setup database: %v", err)
	}
}

func main() {
	// Server Setup
	host, ok := os.LookupEnv("SERVER_HOST")
	if !ok {
		log.Fatalf("SERVER_HOST missing")
	}
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		log.Fatalf("SERVER_PORT missing")
	}
	aud, ok := os.LookupEnv("AUTH0_AUD")
	if !ok {
		log.Fatalf("AUTH0_AUD missing")
	}
	iss, ok := os.LookupEnv("AUTH0_ISS")
	if !ok {
		log.Fatalf("AUTH0_ISS missing")
	}
	corsOrigins, ok := os.LookupEnv("CORS_ORIGINS")
	if !ok {
		log.Fatalf("CORS_ORIGINS missing")
	}
	respAPI, err := server.SetupServer(host, port, iss, aud, corsOrigins)
	if err != nil {
		log.Fatalf("Could not start web server: %v", err)
	}

	pollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL_MINUTES"))
	if err != nil {
		log.Fatalf("could not load poll interval: %v", err)
	}

	// Publisher Setup
	telegramAPIUri, ok := os.LookupEnv("TELEGRAM_CONTACT_URI")
	if !ok {
		log.Fatalf("TELEGRAM_CONTACT_URI missing")
	}
	telegramAuth0Aud, ok := os.LookupEnv("TELEGRAM_AUTH0_AUD")
	if !ok {
		log.Fatalf("TELEGRAM_AUTH0_AUD missing")
	}
	telegramAuth0ClientID, ok := os.LookupEnv("TELEGRAM_AUTH0_CLIENT_ID")
	if !ok {
		log.Fatalf("TELEGRAM_AUTH0_CLIENT_ID missing")
	}
	telegramAuth0ClientSecret, ok := os.LookupEnv("TELEGRAM_AUTH0_CLIENT_SECRET")
	if !ok {
		log.Fatalf("TELEGRAM_AUTH0_CLIENT_SECRET missing")
	}
	telegramAuth0TokenUrl, ok := os.LookupEnv("TELEGRAM_AUTH0_TOKEN_URL")
	if !ok {
		log.Fatalf("TELEGRAM_AUTH0_TOKEN_URL missing")
	}
	tp, err := notifier.NewTelegramPublisher(
		telegramAPIUri,
		telegramAuth0TokenUrl,
		telegramAuth0ClientID,
		telegramAuth0ClientSecret,
		telegramAuth0Aud,
	)
	if err != nil {
		log.Fatalf("Could not setup telegram publisher: %v", err)
	}

	sendGridAPIKey, ok := os.LookupEnv("SENDGRID_API_KEY")
	if !ok {
		log.Fatalf("SENDGRID_API_KEY missing")
	}
	sendGridEmail, ok := os.LookupEnv("SENDGRID_EMAIL")
	if !ok {
		log.Fatalf("SENDGRID_EMAIL missing")
	}
	ep := notifier.NewEmailPublisher(sendGridAPIKey, sendGridEmail)

	c := make(chan model.Covid19Region, 500)
	_ = rki.NewRegionUpdater(time.Duration(pollInterval)*time.Minute, c)
	_ = notifier.NewCovid19Notifier(c, tp, ep)

	// Start Server
	go respAPI.Start()

	done := make(chan os.Signal, 1)
	<-done
}
