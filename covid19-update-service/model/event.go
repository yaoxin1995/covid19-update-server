package model

import (
	"fmt"
	"time"

	"github.com/pmoule/go2hal/hal"

	"github.com/jinzhu/gorm"
)

type EventCollection []Event

// Every time the threshold of a Topic is exceeded an Event is created.
//The Event includes the message which is sent to the notification channels that are defines in the corresponding Subscription.
type Event struct {
	PersistentModel
	Message string `json:"message"`
	TopicID uint   `sql:"type:bigint REFERENCES topics(id) ON DELETE CASCADE" json:"-"`
}

const messagePattern = "The COVID-19 7-day-incidence value at your location (%f, %f) currently is %.2f. You receive this message, because you set the alert threshold to %d. (Update: %s)"

// Creates a new Event for a Topic t using the corresponding Covid19Region c.
func NewEvent(c Covid19Region, t Topic) (Event, error) {
	timestamp := time.Now()
	message := fmt.Sprintf(messagePattern, t.Position.Longitude, t.Position.Latitude, c.Incidence, t.Threshold, timestamp.Format("02.01.2006 15:04 -0700 MST"))
	e := Event{
		Message: message,
		TopicID: t.ID,
	}
	return e, e.Store()
}

// Persists the Event.
func (e *Event) Store() error {
	return db.Save(&e).Error
}

// Loads all Events of the Topic that is identified by tID.
func GetEvents(tID uint) (EventCollection, error) {
	var e EventCollection
	err := db.Where("topic_id = ?", tID).Find(&e).Error
	return e, err
}

// Loads Event with the identifier eID for the Topic tID. If no event is found, nil is returned.
func GetEvent(eID, tID uint) (*Event, error) {
	e := &Event{}
	err := db.Where("id = ? AND topic_id = ?", eID, tID).Find(&e).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

// Loads the newest Events of the Topic with a limit.
func GetEventsWithLimit(tID, limit uint) (EventCollection, error) {
	var e EventCollection
	err := db.Where("topic_id = ?", tID).Limit(limit).Order("created_at DESC").Find(&e).Error
	return e, err
}

// Represents the Event with the JSON Hypertext Application Language.
// path is the relative URI of the Event.
func (e Event) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()
	root.AddData(e)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}

// Represents a collection of Event with the JSON Hypertext Application Language.
// path is the relative URI of the Event collection.
func (ec EventCollection) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	var embeddedEvs []hal.Resource

	for _, e := range ec {
		eHref := fmt.Sprintf("%s/%d", path, e.ID)
		eSelfLink, _ := hal.NewLinkObject(eHref)

		eSelfRel, _ := hal.NewLinkRelation("self")
		eSelfRel.SetLink(eSelfLink)

		embeddedEv := hal.NewResourceObject()
		embeddedEv.AddLink(eSelfRel)
		embeddedEv.AddData(e)
		embeddedEvs = append(embeddedEvs, embeddedEv)
	}

	evs, _ := hal.NewResourceRelation("events")
	evs.SetResources(embeddedEvs)
	root.AddResource(evs)

	return root
}
