package model

import (
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	CommonModelFields
	Email          *string `json:"email"`
	TelegramChatID *string `json:"telegramChatId"`
	Topics         []Topic `json:"-"`
}

func NewSubscription(email, telegram *string) (Subscription, error) {
	s := Subscription{
		Email:          email,
		TelegramChatID: telegram,
		Topics:         []Topic{},
	}
	err := s.Store()
	return s, err
}

func (s *Subscription) Store() error {
	return db.Save(&s).Error
}

func GetSubscription(id uint) (*Subscription, error) {
	s := &Subscription{}
	err := db.First(s, id).Preload("Topics").Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func GetSubscriptions() ([]Subscription, error) {
	var subs []Subscription
	err := db.Find(&subs).Preload("Topics").Error
	return subs, err
}

func (s *Subscription) Delete() error {
	return db.Unscoped().Delete(s).Error
}

func (s *Subscription) Update(email, telegram *string) error {
	s.Email = email
	s.TelegramChatID = telegram
	return s.Store()
}
