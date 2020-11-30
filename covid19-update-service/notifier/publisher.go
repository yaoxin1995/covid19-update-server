package notifier

import "covid19-update-service/model"

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
	// ToDo Publish via Telegram
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
	return nil
}
