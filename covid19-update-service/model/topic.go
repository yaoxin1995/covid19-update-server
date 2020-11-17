package model

type Topic struct {
	CommonModelFields
	Position       GPSPosition `gorm:"embedded;embeddedPrefix:position_" ,json:"position"`
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

/*func GetTopic(tID, sID uint) (Topic, error) {
	t := Topic{}
	err := db.Where("id = ? AND subscription_id = ?", tID, sID).Error
	if err != nil {

	}
}*/

func (t *Topic) Delete() error {
	return db.Where("id = ?", t.ID).Delete(&Topic{}).Error
}
