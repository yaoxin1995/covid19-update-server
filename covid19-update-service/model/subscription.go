package model

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

type Subscription struct {
	CommonModelFields
	Email          null.String `json:"email"`
	TelegramChatID null.String `json:"telegramChatId"`
	Topics         []Topic     `json:"-"`
}

func NewSubscription(email, telegram *string) (Subscription, error) {
	s := Subscription{
		Email:          null.StringFromPtr(email),
		TelegramChatID: null.StringFromPtr(telegram),
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
	err := db.Preload("Topics").First(s, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func SubscriptionExists(id uint) bool {
	s := &Subscription{}
	err := db.First(s, id).Error
	if err != nil {
		return false
	}
	return true
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
	s.Email = null.StringFromPtr(email)
	s.TelegramChatID = null.StringFromPtr(telegram)
	return s.Store()
}
