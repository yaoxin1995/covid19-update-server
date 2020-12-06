package notifier

import (
	"bytes"
	"covid19-update-service/model"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type Publisher interface {
	Publish(address string, e model.Event) error
}

// Telegram

var TelegramApiURI = ""

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
	// ToDo Publish via email
	log.Printf("Sent Email")
	return nil
}
