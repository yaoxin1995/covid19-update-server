package rki

import (
	"covid19-update-service/model"
	"log"
	"time"
)

// Poller with async polling
// Later it should then trigger notification checker (internal)

type Covid19Poller struct {
	ticker *time.Ticker
}

func NewCovid19Poller(interval time.Duration) *Covid19Poller {
	cp := &Covid19Poller{
		ticker: time.NewTicker(interval),
	}
	go cp.run()
	return cp
}

func (cp *Covid19Poller) run() error {
	for {
		select {
		case <-cp.ticker.C:
			log.Printf("Start polling...")
			successfulPolls, totalPolls, err := poll()
			if err != nil {
				log.Printf("Could not poll incidences: %v", err)
				break
			}
			log.Printf("Polling finished successfully. %d of %d incidences were updated.", successfulPolls, totalPolls)
			// ToDo Trigger Notification Logic
		}
	}
}

func poll() (int, int, error) {
	rsp, err := GetAllIncidences()
	if err != nil {
		return 0, 0, err
	}
	updateSuccessCnt := 0
	for _, f := range rsp.Features {
		_, err := model.NewIncidence(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err == nil {
			updateSuccessCnt++
		}
	}
	return updateSuccessCnt, len(rsp.Features), nil
}
