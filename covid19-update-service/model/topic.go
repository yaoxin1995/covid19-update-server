package model

import "github.com/jinzhu/gorm"

type Topic struct {
	CommonModelFields
	Position       GPSPosition `gorm:"embedded;embeddedPrefix:position_" json:"position"`
	Threshold      uint        `json:"threshold"`
	SubscriptionID uint        `sql:"type:bigint REFERENCES subscriptions(id) ON DELETE CASCADE" json:"-"`
}

type GPSPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewTopic(position GPSPosition, threshold, subID uint) (Topic, error) {
	t := Topic{
		Position:       position,
		Threshold:      threshold,
		SubscriptionID: subID,
	}
	err := t.Store()
	return t, err
}

func (t *Topic) Store() error {
	return db.Save(&t).Error
}

func (t *Topic) Update(position GPSPosition, threshold uint) error {
	t.Position = position
	t.Threshold = threshold
	return t.Store()
}

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

func GetTopics(sID uint) ([]Topic, error) {
	var tops []Topic
	err := db.Where("subscription_id = ?", sID).Find(&tops).Error
	return tops, err
}

func (t *Topic) Delete() error {
	return db.Unscoped().Delete(t).Error
}
