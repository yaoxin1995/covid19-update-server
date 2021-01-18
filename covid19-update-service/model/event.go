package model

import (
	"fmt"

	"github.com/pmoule/go2hal/hal"

	"github.com/jinzhu/gorm"
)

type EventCollection []Event

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

func GetEvents(tID uint) (EventCollection, error) {
	var e EventCollection
	err := db.Where("topic_id = ?", tID).Find(&e).Error
	return e, err
}

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

func GetEventsWithLimit(tID, limit uint) (EventCollection, error) {
	var e EventCollection
	err := db.Where("topic_id = ?", tID).Limit(limit).Order("created_at DESC").Find(&e).Error
	return e, err
}

func (e Event) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()
	root.AddData(e)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return root
}

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
