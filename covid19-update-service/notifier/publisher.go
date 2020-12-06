package notifier

import (
	"bytes"
	"covid19-update-service/model"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Publisher interface {
	Publish(address string, e model.Event) error
}

// Telegram

var TelegramApiURI = ""
var SendGridAPIKey = ""

type TelegramPublisher struct {
	ChatID string
}

func NewTelegramPublisher(chatID string) TelegramPublisher {
	return TelegramPublisher{chatID}
}

func (tp *TelegramPublisher) Publish(e model.Event) error {
	values := map[string]io.Reader{
		"recipient": strings.NewReader(tp.ChatID), // lets assume its this file
		"msg":       strings.NewReader(e.Message),
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		fw, _ = w.CreateFormField(key)
		_, _ = io.Copy(fw, r)
	}
	_ = w.Close()
	req, err := http.NewRequest("POST", TelegramApiURI, &b)
	if err != nil {
		return fmt.Errorf("could not create telegram request: %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send telegram request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// Email

type EmailPublisher struct {
	Email string
}

func NewEmailPublisher(email string) EmailPublisher {
	return EmailPublisher{email}
}

func (ep *EmailPublisher) Publish(e model.Event) error {
	from := mail.NewEmail("Covid 19 Updater", "ludwig_maximilian.leibl@mailbox.tu-dresden.de")
	subject := "Covid19 Update"
	to := mail.NewEmail(ep.Email, ep.Email)
	message := mail.NewSingleEmail(from, subject, to, e.Message, e.Message)
	client := sendgrid.NewSendClient(SendGridAPIKey)
	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("could not send email request: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Body)
	}
	return nil
}
