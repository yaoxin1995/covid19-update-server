package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pmoule/go2hal/hal"
	"gopkg.in/guregu/null.v3"
)

type SubscriptionCollection []Subscription

type Subscription struct {
	CommonModelFields
	Email          null.String     `json:"email"`
	TelegramChatID null.String     `json:"telegramChatId"`
	Topics         TopicCollection `json:"-"`
	ClientID       string          `json:"-"`
}

func NewSubscription(email, telegram *string, clientID string, topics TopicCollection) (Subscription, error) {
	s := Subscription{
		Email:          null.StringFromPtr(email),
		TelegramChatID: null.StringFromPtr(telegram),
		Topics:         topics,
		ClientID:       clientID,
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

func GetSubscriptions(owner string) (SubscriptionCollection, error) {
	var subs SubscriptionCollection
	err := db.Where("client_id=?", owner).Find(&subs).Error
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

func (s Subscription) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()

	// Add email and telegram manually, because HAL JSON encoder does not can handle null.String
	data := root.Data()
	data["id"] = s.ID
	data["email"] = s.Email.Ptr()
	data["telegramChatId"] = s.TelegramChatID.Ptr()
	root.AddData(data)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	topicsRel, _ := hal.NewLinkRelation("topics")
	topicsLink := &hal.LinkObject{Href: fmt.Sprintf("%s/topics", path)}
	topicsRel.SetLink(topicsLink)
	root.AddLink(topicsRel)

	return root
}

func (sc SubscriptionCollection) ToHAL(path string) hal.Resource {
	root := hal.NewResourceObject()

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: path}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	var embeddedSubs []hal.Resource

	for _, s := range sc {
		eHref := fmt.Sprintf("%s/%d", path, s.ID)
		eSelfLink, _ := hal.NewLinkObject(eHref)

		eSelfRel, _ := hal.NewLinkRelation("self")
		eSelfRel.SetLink(eSelfLink)

		embeddedSub := hal.NewResourceObject()
		embeddedSub.AddLink(eSelfRel)
		data := embeddedSub.Data()
		data["id"] = s.ID
		data["email"] = s.Email.Ptr()
		data["telegramChatId"] = s.TelegramChatID.Ptr()
		embeddedSub.AddData(data)
		embeddedSubs = append(embeddedSubs, embeddedSub)
	}

	subs, _ := hal.NewResourceRelation("subscriptions")
	subs.SetResources(embeddedSubs)
	root.AddResource(subs)

	return root
}
