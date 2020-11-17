package model

type Subscription struct {
	CommonModelFields
	Email    *string `json:"email"`
	Telegram *string `json:"telegram"`
	Topics   []Topic `gorm:"preload:true" json:"-"`
}

func NewSubscription(email, telegram *string) (Subscription, error) {
	s := Subscription{
		Email:    email,
		Telegram: telegram,
		Topics:   []Topic{},
	}
	err := s.Store()
	return s, err
}

func (s *Subscription) Store() error {
	return db.Save(&s).Error
}

func GetSubscription(id uint) (Subscription, error) {
	s := Subscription{}
	err := db.First(&s, id).Error
	return s, err
}

func GetSubscriptions() ([]Subscription, error) {
	var subs []Subscription
	err := db.Find(&subs).Error
	return subs, err
}

func DeleteSubscription(id uint) error {
	return db.Where("id = ?", id).Delete(&Subscription{}).Error
}

func (s *Subscription) Update(email, telegram *string) error {
	s.Email = email
	s.Telegram = telegram
	return s.Store()
}
