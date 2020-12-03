package notifier

import (
	"covid19-update-service/model"
	"log"
)

type Publisher interface {
	Publish(e model.Event) error
}

type TelegramPublisher struct {
	ChatID string
}

func NewTelegramPublisher(chatID string) TelegramPublisher {
	return TelegramPublisher{chatID}
}

func (tp *TelegramPublisher) Publish(e model.Event) error {
	// ToDo Publish via TelegramChatID
	log.Printf("Sent telegram")
	return nil
}

type EmailPublisher struct {
	Email string
}

func NewEmailPublisher(email string) EmailPublisher {
	return EmailPublisher{email}
}

func (ep *EmailPublisher) Publish(e model.Event) error {
	// ToDo Publish via email
	log.Printf("Sent Email")
	return nil
}
