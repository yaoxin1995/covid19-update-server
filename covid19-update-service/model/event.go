package model

import "fmt"

type Event struct {
	CommonModelFields
	Message string `json:"message"`
	TopicID uint   `sql:"type:bigint REFERENCES topics(id) ON DELETE CASCADE" json:"-"`
}

const messagePattern = "The Covid 19 incidence value at your location (%f, %f) currently is %f. You receive this message, because you set the threshold to %d."

func NewEvent(i Incidence, t Topic) (Event, error) {
	message := fmt.Sprintf(messagePattern, t.Position.Longitude, t.Position.Latitude, i.Cases7Per100k, t.Threshold)
	e := Event{
		Message: message,
		TopicID: t.ID,
	}
	return e, e.Store()
}

func (e *Event) Store() error {
	return db.Save(&e).Error
}
