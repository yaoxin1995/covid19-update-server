package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Handles the OAuth 2.0 client credentials flow with https://auth0.com and automatically refreshes the access token.
type Auth0AccessTokenHelper struct {
	mu                sync.RWMutex
	ticker            *time.Ticker
	accessToken       string
	OAuthTokenUrl     string
	OAuthClientID     string
	OAuthClientSecret string
	OAuthAudience     string
}

type auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Creates new Auth0AccessTokenHelper given the Auth0 token URL, client ID, client secret and audience.
func NewAuth0AccessTokenHelper(tokenUrl, cID, cSecret, aud string) (*Auth0AccessTokenHelper, error) {
	a := &Auth0AccessTokenHelper{
		mu:                sync.RWMutex{},
		OAuthTokenUrl:     tokenUrl,
		OAuthClientID:     cID,
		OAuthClientSecret: cSecret,
		OAuthAudience:     aud,
	}
	tokenResponse, err := a.requestAccessToken()
	if err != nil {
		return nil, fmt.Errorf("could not get initial access token: %v", err)
	}
	a.accessToken = tokenResponse.AccessToken
	a.ticker = time.NewTicker(time.Duration(tokenResponse.ExpiresIn)*time.Second - time.Hour)
	go a.scheduleTokenRefresh()
	return a, nil
}

func (a *Auth0AccessTokenHelper) scheduleTokenRefresh() {
	for {
		select {
		case <-a.ticker.C:
			tokenResponse, err := a.requestAccessToken()
			if err != nil {
				log.Printf("Auth0 token refresh failed: %v", err)
				return
			}
			a.mu.Lock()
			a.accessToken = tokenResponse.AccessToken
			a.mu.Unlock()
			a.ticker.Stop()
			a.ticker = time.NewTicker(time.Duration(tokenResponse.ExpiresIn)*time.Second - time.Hour)
			log.Printf("Refreshed Auth0 access token.")
		}
	}
}

// Returns valid access token from Auth0.
func (a *Auth0AccessTokenHelper) GetAccessToken() string {
	a.mu.RLock()
	token := a.accessToken
	defer a.mu.RUnlock()
	return token
}

func (a *Auth0AccessTokenHelper) requestAccessToken() (*auth0TokenResponse, error) {
	payload := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&audience=%s", a.OAuthClientID, a.OAuthClientSecret, a.OAuthAudience))

	req, err := http.NewRequest("POST", a.OAuthTokenUrl, payload)
	if err != nil {
		return nil, fmt.Errorf("could not create access token request: %v", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send access token request: %v", err)
	}

	var tokenResponse auth0TokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("could not decode access token response: %v", err)
	}
	return &tokenResponse, nil
}
