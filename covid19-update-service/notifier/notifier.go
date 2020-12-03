package notifier

import (
	"covid19-update-service/model"
	"fmt"
	"log"
)

type Covid19Notifier struct {
	c <-chan model.Covid19Region
}

func NewCovid19Notifier(c <-chan model.Covid19Region) *Covid19Notifier {
	cn := &Covid19Notifier{
		c: c,
	}
	go cn.run()
	return cn
}

func (cn *Covid19Notifier) run() error {
	for {
		select {
		case cov19region := <-cn.c:
			log.Printf("Creating notifications for region %d...", cov19region.ID)
			notify(cov19region)
		}
	}
}

func notify(cov19region model.Covid19Region) {
	tops, err := model.GetTopicsWithThresholdAlert(cov19region)
	if err != nil {
		fmt.Printf("Could not load topics for notification region: %v", err)
	}
	for _, t := range tops {
		e, err := model.NewEvent(cov19region, t)
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

func shipEvent(e model.Event, sID uint) error {
	s, err := model.GetSubscription(sID)
	if err != nil {
		return fmt.Errorf("could not load subscription: %v", err)
	}
	if s.TelegramChatID != nil {
		tp := NewTelegramPublisher(*s.TelegramChatID)
		err = tp.Publish(e)
		if err != nil {
			fmt.Printf("Could not publish event via telegram: %v", err)
		}
	}
	if s.Email != nil {
		ep := NewEmailPublisher(*s.TelegramChatID)
		err = ep.Publish(e)
		if err != nil {
			fmt.Printf("Could not publish event via email: %v", err)
		}
	}
	return nil
}
