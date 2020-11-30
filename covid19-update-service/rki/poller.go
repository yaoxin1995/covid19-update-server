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
			log.Printf("Updating Covid 19 incidence values...")
			updatedIncidenceIDs, err := updateIncidences()
			if err != nil {
				log.Printf("Could perform incidence update: %v", err)
				break
			}
			log.Printf("Updating Covid 19 incidence values finished successfully. %d regions were updated.", len(*updatedIncidenceIDs))
			// ToDo Find topics that are associated with updated value
		}
	}
}

func updateIncidences() (*[]uint, error) {
	rsp, err := GetAllIncidences()
	if err != nil {
		return nil, err
	}
	updatedIncidenceIDs := make([]uint, 0)
	for _, f := range rsp.Features {
		isNew, err := updateIfNew(&f)
		if err != nil {
			log.Printf("Could not update incidence: %v", err)
		}
		if isNew {
			updatedIncidenceIDs = append(updatedIncidenceIDs, f.Attributes.ObjectID)
		}
	}
	return &updatedIncidenceIDs, err
}

func updateIfNew(f *Feature) (bool, error) {
	i, err := model.GetIncidence(f.Attributes.ObjectID)
	if err != nil {
		return false, err
	}
	// Entry does not exist
	if i == nil {
		_, err := model.NewIncidence(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Entry does exist
	newTime, err := f.Attributes.LastUpdate()
	if err != nil {
		return false, nil
	}
	if i.UpdatedAt.UTC().Before(newTime.UTC()) {
		_, err := model.NewIncidence(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
