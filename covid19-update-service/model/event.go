package model

import "fmt"

type Event struct {
	CommonModelFields
	Message string `json:"message"`
	TopicID uint   `sql:"type:bigint REFERENCES topics(id) ON DELETE CASCADE" json:"-"`
}

const messagePattern = "The Covid 19 7-day-incidence value at your location (%f, %f) currently is %f. You receive this message, because you set the alert threshold to %d."

func NewEvent(c Covid19Region, t Topic) (Event, error) {
	message := fmt.Sprintf(messagePattern, t.Position.Longitude, t.Position.Latitude, c.Incidence, t.Threshold)
	e := Event{
		Message: message,
		TopicID: t.ID,
	}
	return e, e.Store()
}

func (e *Event) Store() error {
	return db.Save(&e).Error
}
