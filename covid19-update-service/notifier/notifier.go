package notifier

import (
	"covid19-update-service/model"
	"fmt"
	"log"
)

type Covid19Notifier struct {
	c  <-chan model.Covid19Region
	tp TelegramPublisher
	ep EmailPublisher
}

func NewCovid19Notifier(c <-chan model.Covid19Region, tPub TelegramPublisher, ePub EmailPublisher) *Covid19Notifier {
	cn := &Covid19Notifier{
		c:  c,
		tp: tPub,
		ep: ePub,
	}
	go cn.run()
	return cn
}

func (cn *Covid19Notifier) run() error {
	for {
		select {
		case cov19region := <-cn.c:
			// log.Printf("Creating notifications for region %d...", cov19region.ID)
			cn.notify(cov19region)
		}
	}
}

func (cn *Covid19Notifier) notify(cov19region model.Covid19Region) {
	tops, err := model.GetTopicsWithThresholdAlert(cov19region)
	if err != nil {
		log.Printf("Could not load topics for notification region: %v", err)
	}
	for _, t := range tops {
		e, err := model.NewEvent(cov19region, t)
		if err != nil {
			log.Printf("Could not create event: %v", err)
			continue
		}
		err = cn.shipEvent(e, t.SubscriptionID)
		if err != nil {
			log.Printf("Could not ship event: %v", err)
		}
	}
}

func (cn *Covid19Notifier) shipEvent(e model.Event, sID uint) error {
	s, err := model.GetSubscription(sID)
	if err != nil {
		return fmt.Errorf("could not load subscription: %v", err)
	}
	if s.TelegramChatID.Ptr() != nil {
		go func() {
			log.Printf("Sending telegram notification...")
			err = cn.tp.Publish(s.TelegramChatID.String, e)
			if err != nil {
				log.Printf("Could not publish event %d via telegram: %v", e.ID, err)
				return
			}
			log.Printf("Sending telegram notification finished successfully")
		}()
	}
	if s.Email.Ptr() != nil {
		go func() {
			log.Printf("Sending email notification...")
			err = cn.ep.Publish(s.Email.String, e)
			if err != nil {
				log.Printf("Could not publish event %d via telegram: %v", e.ID, err)
				return
			}
			log.Printf("Sending email notification finished successfully")
		}()
	}
	return nil
}
