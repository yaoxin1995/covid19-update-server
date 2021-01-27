package model

import (
	"encoding/json"
	"fmt"

	"github.com/pmoule/go2hal/hal"

	"github.com/jinzhu/gorm"
)

type TopicCollection []Topic

// A Topic represents a location that should be monitored by the service.
type Topic struct {
	PersistentModel
	Position        GPSPosition     `gorm:"embedded;embeddedPrefix:position_" json:"position"`
	Threshold       uint            `json:"threshold"`
	SubscriptionID  uint            `sql:"type:bigint REFERENCES subscriptions(id) ON DELETE CASCADE" json:"-"`
	Covid19RegionID uint            `json:"-"`
	Events          EventCollection `json:"-"`
}

// GPSPosition represents a GPS location.
type GPSPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Unmarshaler for GPSPosition which makes all field requried.
func (p *GPSPosition) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}{}
	all := struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Latitude == nil {
		err = fmt.Errorf("latitude is missing")
	} else if required.Longitude == nil {
		err = fmt.Errorf("longitude is missing")
	} else {
		err = json.Unmarshal(data, &all)
		p.Longitude = all.Longitude
		p.Latitude = all.Latitude
	}
	return err
}

// Creates a new Topic.
// The GPSPosition hat to be provided with the ID of the associated Covid19Region cov19RegID.
func NewTopic(position GPSPosition, threshold, subID uint, cov19RegID uint) (Topic, error) {
	t := Topic{
		Position:        position,
		Threshold:       threshold,
		SubscriptionID:  subID,
		Covid19RegionID: cov19RegID,
	}
	err := t.Store()
	return t, err
}

// Persists a Topic.
func (t *Topic) Store() error {
	return db.Save(&t).Error
}

// Updates Topic.
// The GPSPosition hat to be provided with the ID of the associated Covid19Region cov19RegID.
func (t *Topic) Update(position GPSPosition, threshold uint, cov19RegID uint) error {
	t.Position = position
	t.Threshold = threshold
	t.Covid19RegionID = cov19RegID
	return t.Store()
}

// Returns a Topic with the given tID scoped to the Subscription sID or nil if no such Topic was found.
func GetTopic(tID, sID uint) (*Topic, error) {
	t := &Topic{}
	err := db.Where("id = ? AND subscription_id = ?", tID, sID).First(t).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

// Returns all Topics for the Subscription sID.
func GetTopicsBySubscriptionID(sID uint) (TopicCollection, error) {
	var tops TopicCollection
	err := db.Where("subscription_id = ?", sID).Find(&tops).Error
	return tops, err
}

// Returns all Topics for a CovidRegions whose threshold is below the regions incidence value.
func GetTopicsWithThresholdAlert(c Covid19Region) (TopicCollection, error) {
	var tops TopicCollection
	err := db.Where("covid19_region_id = ? AND threshold <= ?", c.ID, c.Incidence).Find(&tops).Error
	return tops, err
}

// Deletes a Topic from the database.
func (t *Topic) Delete() error {
	return db.Unscoped().Delete(t).Error
}

// Represents a collection of TopicCollection with the JSON Hypertext Application Language.
// path is the relative URI of the TopicCollection collection.
func (tc TopicCollection) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	var embeddedTops []hal.Resource

	for _, t := range tc {
		eHref := fmt.Sprintf("%s/%d", path, t.ID)
		eSelfLink, _ := hal.NewLinkObject(eHref)

		eSelfRel, _ := hal.NewLinkRelation("self")
		eSelfRel.SetLink(eSelfLink)

		embeddedTop := hal.NewResourceObject()
		embeddedTop.AddLink(eSelfRel)
		embeddedTop.AddData(t)
		embeddedTops = append(embeddedTops, embeddedTop)
	}

	tops, _ := hal.NewResourceRelation("topics")
	tops.SetResources(embeddedTops)
	root.AddResource(tops)

	return root
}

// Represents the Topic with the JSON Hypertext Application Language.
// path is the relative URI of the Topic.
func (t Topic) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()
	root.AddData(t)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	eventsRel, _ := hal.NewLinkRelation("events")
	eventLink := &hal.LinkObject{Href: fmt.Sprintf("%s/events", path)}
	eventsRel.SetLink(eventLink)
	root.AddLink(eventsRel)

	incidenceRel, _ := hal.NewLinkRelation("incidence")
	incidenceLink := &hal.LinkObject{Href: fmt.Sprintf("%s/incidence", path)}
	incidenceRel.SetLink(incidenceLink)
	root.AddLink(incidenceRel)

	return root
}
