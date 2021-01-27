package notifier

import (
	"covid19-update-service/model"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Telegram

type TelegramPublisher struct {
	TelegramServiceURI string
	accessTokenHelper  *Auth0AccessTokenHelper
}

// Creates new TelegramPublisher for the Telegram Notification Service.
func NewTelegramPublisher(tServiceUri string, auth0Helper *Auth0AccessTokenHelper) *TelegramPublisher {
	tp := &TelegramPublisher{
		TelegramServiceURI: tServiceUri,
		accessTokenHelper:  auth0Helper,
	}
	return tp
}

// Publishes the Event e to the chat identified by chatID via the Telegram Notification Service.
func (tp *TelegramPublisher) Publish(chatID string, e model.Event) error {
	data := url.Values{}
	data.Set("recipient", chatID)
	data.Set("msg", e.Message)

	req, err := http.NewRequest("POST", tp.TelegramServiceURI, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("could not create telegram request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tp.accessTokenHelper.GetAccessToken()))
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
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
	SendGridAPIKey string
	SendGridEmail  string
}

// Creates new EmailPublisher the SendGrid credential and a valid email, that will be used as sender address.
func NewEmailPublisher(sendGridApiKey, sendGridEmail string) *EmailPublisher {
	return &EmailPublisher{SendGridAPIKey: sendGridApiKey, SendGridEmail: sendGridEmail}
}

// Publishes the Event e to the email address using the SendGrid Email Delivery Service.
func (ep *EmailPublisher) Publish(email string, e model.Event) error {
	from := mail.NewEmail("Covid 19 Updater", ep.SendGridEmail)
	subject := "Covid19 Update"
	to := mail.NewEmail(email, email)
	message := mail.NewSingleEmail(from, subject, to, e.Message, e.Message)
	client := sendgrid.NewSendClient(ep.SendGridAPIKey)
	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("could not send email request: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Body)
	}
	return nil
}
