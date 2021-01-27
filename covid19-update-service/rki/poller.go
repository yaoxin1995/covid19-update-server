package rki

import (
	"covid19-update-service/model"
	"log"
	"time"
)

// Regularly polls the current incidence values from the RKI Corona Landkreise REST API and updates the model.Covid19Region if necessary.
type Covid19RegionUpdater struct {
	ticker *time.Ticker
	c      chan<- model.Covid19Region
}

// Creates new Covid19RegionUpdater which publishes the updated model.Covid19Region to the buffered channel c.
func NewRegionUpdater(interval time.Duration, c chan<- model.Covid19Region) *Covid19RegionUpdater {
	cp := &Covid19RegionUpdater{
		ticker: time.NewTicker(interval),
		c:      c,
	}
	go cp.run()
	return cp
}

func (cp *Covid19RegionUpdater) run() error {
	for {
		select {
		case <-cp.ticker.C:
			log.Printf("Updating Covid 19 regions...")
			updatedRegions, err := updateRegions()
			if err != nil {
				log.Printf("Could not perform update: %v", err)
				break
			}
			log.Printf("Updating Covid 19 regions finished successfully. %d regions were updated.", len(*updatedRegions))
			for _, r := range *updatedRegions {
				cp.c <- r
			}
		}
	}
}

func updateRegions() (*[]model.Covid19Region, error) {
	rsp, err := getAllRegions()
	if err != nil {
		return nil, err
	}
	updatedRegionsIDs := make([]model.Covid19Region, 0)
	for _, f := range rsp.Features {
		i, err := updateIfNew(&f)
		if err != nil {
			log.Printf("Could not update region: %v", err)
		}
		if i != nil {
			updatedRegionsIDs = append(updatedRegionsIDs, *i)
		}
	}
	return &updatedRegionsIDs, err
}

func updateIfNew(f *feature) (*model.Covid19Region, error) {
	i, err := model.GetCovid19Region(f.Attributes.ObjectID)
	if err != nil {
		return nil, err
	}
	// Entry does not exist
	if i == nil {
		inc, err := model.NewCovid19Region(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err != nil {
			return nil, err
		}
		return &inc, nil
	}

	// Entry does exist
	newTime, err := f.Attributes.lastUpdate()
	if err != nil {
		return nil, nil
	}
	if i.UpdatedAt.UTC().Before(newTime.UTC()) { // Check if existing entry is outdated
		inc, err := model.NewCovid19Region(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err != nil {
			return nil, err
		}
		return &inc, nil
	}
	return nil, nil
}
