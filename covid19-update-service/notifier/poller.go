package notifier

import (
	"covid19-update-service/model"
	"covid19-update-service/rki"
	"fmt"
	"log"
	"time"
)

type Covid19Notifier struct {
	ticker *time.Ticker
}

func NewCovid19Notifier(interval time.Duration) *Covid19Notifier {
	cp := &Covid19Notifier{
		ticker: time.NewTicker(interval),
	}
	go cp.run()
	return cp
}

func (cp *Covid19Notifier) run() error {
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
			notify(updatedIncidenceIDs)
		}
	}
}

func notify(incs *[]*model.Incidence) {
	for _, i := range *incs {
		tops, err := model.GetTopicsByIncidence(*i)
		if err != nil {
			log.Printf("Could not load topics for incidence: %v", err)
			continue
		}
		for _, t := range tops {
			e, err := model.NewEvent(*i, t)
			if err != nil {
				log.Printf("Could not create event: %v", err)
				continue
			}
			err = shipEvent(e, t.SubscriptionID)
			if err != nil {
				log.Printf("Could not ship event: %v", err)
			}
		}
	}
}

func shipEvent(e model.Event, sID uint) error {
	s, err := model.GetSubscription(sID)
	if err != nil {
		return fmt.Errorf("could not load subscription: %v", err)
	}
	if s.Telegram != nil {
		tp := NewTelegramPublisher(*s.Telegram)
		err = tp.Publish(e)
		if err != nil {
			return fmt.Errorf("could not publish via telegram: %v", err)
		}
	}
	if s.Email != nil {
		ep := NewEmailPublisher(*s.Telegram)
		err = ep.Publish(e)
		if err != nil {
			return fmt.Errorf("could not publish via email: %v", err)
		}
	}
	return nil
}

func updateIncidences() (*[]*model.Incidence, error) {
	rsp, err := rki.GetAllIncidences()
	if err != nil {
		return nil, err
	}
	updatedIncidenceIDs := make([]*model.Incidence, 0)
	for _, f := range rsp.Features {
		i, err := updateIfNew(&f)
		if err != nil {
			log.Printf("Could not update incidence: %v", err)
		}
		if i != nil {
			updatedIncidenceIDs = append(updatedIncidenceIDs, i)
		}
	}
	return &updatedIncidenceIDs, err
}

func updateIfNew(f *rki.Feature) (*model.Incidence, error) {
	i, err := model.GetIncidence(f.Attributes.ObjectID)
	if err != nil {
		return nil, err
	}
	// Entry does not exist
	if i == nil {
		inc, err := model.NewIncidence(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err != nil {
			return nil, err
		}
		return &inc, nil
	}

	// Entry does exist
	newTime, err := f.Attributes.LastUpdate()
	if err != nil {
		return nil, nil
	}
	if i.UpdatedAt.UTC().Before(newTime.UTC()) {
		inc, err := model.NewIncidence(f.Attributes.ObjectID, f.Attributes.Cases7Per100k)
		if err != nil {
			return nil, err
		}
		return &inc, nil
	}
	return nil, nil
}
