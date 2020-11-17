package model

import (
	"github.com/jinzhu/gorm"
)

type Subscription struct {
	CommonModelFields
	Email    *string `json:"email"`
	Telegram *string `json:"telegram"`
	Topics   []Topic `json:"-"`
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

// ToDo Pointer nutzen: nil wenn nicht gefunden error für alles andere: in request handler auf nil überprüfen!
// https://stackoverflow.com/questions/55372748/proper-error-handling-when-no-entity-is-found

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
	return db.Where("id = ?", s.ID).Preload("Topics").Delete(&Subscription{}).Error
}

func (s *Subscription) Update(email, telegram *string) error {
	s.Email = email
	s.Telegram = telegram
	return s.Store()
}

/*
func (s *Subscription) AddTopic(position GPSPosition, threshold uint) (Topic, error) {
	t, err := newTopic(position, threshold, s.ID)
	if err != nil {
		return Topic{}, err
	}
	s.Topics = append(s.Topics, t)
	return t, s.Store()
}

func (s *Subscription) GetTopic(tID uint) (Topic, error) {
	for _, t := range s.Topics {
		if t.ID == tID {
			return t, nil
		}
	}
	return Topic{}, errors.New("subscription does not manage requested topic")
}

func (s *Subscription) GetTopics() []Topic {
	return s.Topics
}

func (s *Subscription) DeleteTopic(t *Topic) error {
	rmI := -1
	for i, tp := range s.Topics {
		if tp.ID == t.ID {
			rmI = i
			break
		}
	}
	if rmI < 0 {
		return errors.New("subscription does not mange requested topic")
	}
	err := t.delete()
	if err != nil {
		return err
	}
	s.Topics = append(s.Topics[:rmI], s.Topics[rmI+1:]...)
	return s.Store()
}
*/
