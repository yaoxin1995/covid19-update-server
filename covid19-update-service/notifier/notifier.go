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
			// log.Printf("Creating notifications for region %d...", cov19region.ID)
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
	if s.TelegramChatID.Ptr() != nil {
		go func() {
			log.Printf("Sending telegram notification...")
			tp := NewTelegramPublisher(s.TelegramChatID.String)
			err = tp.Publish(e)
			if err != nil {
				fmt.Printf("Could not publish event %d via telegram: %v", e.ID, err)
				return
			}
			log.Printf("Sending telegram notification finished successfully")
		}()
	}
	if s.Email.Ptr() != nil {
		go func() {
			log.Printf("Sending email notification...")
			ep := NewEmailPublisher(s.Email.String)
			err = ep.Publish(e)
			if err != nil {
				fmt.Printf("Could not publish event %d via telegram: %v", e.ID, err)
				return
			}
			log.Printf("Sending email notification finished successfully")
		}()
	}
	return nil
}
