package notifier

import (
	"covid19-update-service/model"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Telegram

type TelegramPublisher struct {
	TelegramServiceURI string
	OAuthTokenUrl      string
	OAuthClientID      string
	OAuthClientSecret  string
	OAuthAudience      string
	accessToken        string
}

func NewTelegramPublisher(tServiceUri, tokenUrl, cID, cSecret, aud string) (TelegramPublisher, error) {
	tp := TelegramPublisher{
		TelegramServiceURI: tServiceUri,
		OAuthTokenUrl:      tokenUrl,
		OAuthClientID:      cID,
		OAuthClientSecret:  cSecret,
		OAuthAudience:      aud,
	}
	err := tp.getAccessToken()
	if err != nil {
		return TelegramPublisher{}, fmt.Errorf("could not get initial access token: %v", err)
	}
	return tp, nil
}

func (tp *TelegramPublisher) scheduleAccessTokenRefresh(d time.Duration) {
	time.Sleep(d)
	err := tp.getAccessToken()
	log.Printf("could not refresh telegram access token: %v", err)
}

func (tp *TelegramPublisher) getAccessToken() error {
	type Auth0TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	payload := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&audience=%s", tp.OAuthClientID, tp.OAuthClientSecret, tp.OAuthAudience))

	req, err := http.NewRequest("POST", tp.OAuthTokenUrl, payload)
	if err != nil {
		return fmt.Errorf("could not create access token request: %v", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send access token request: %v", err)
	}

	var tokenResponse Auth0TokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return fmt.Errorf("could not decode access token response: %v", err)
	}

	tp.accessToken = tokenResponse.AccessToken
	go tp.scheduleAccessTokenRefresh(time.Duration(tokenResponse.ExpiresIn)*time.Second - time.Hour)
	return nil
}

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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tp.accessToken))
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

func NewEmailPublisher(sendGridApiKey, sendGridEmail string) EmailPublisher {
	return EmailPublisher{SendGridAPIKey: sendGridApiKey, SendGridEmail: sendGridEmail}
}

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
