package twitch

import (
	"encoding/json"
	"net/http"
	"time"
)

// User twitch user
type User struct {
	ID            string    `json:"_id"`
	Bio           string    `json:"bio"`
	CreatedAt     time.Time `json:"created_at"`
	DisplayName   string    `json:"display_name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Logo          string    `json:"logo"`
	Name          string    `json:"name"`
	Notifications struct {
		Email bool `json:"email"`
		Push  bool `json:"push"`
	} `json:"notifications"`
	Partnered        bool      `json:"partnered"`
	TwitterConnected bool      `json:"twitter_connected"`
	Type             string    `json:"type"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// API twitch resp api
type API struct {
	clientID string
}

// NewAPI creates API
func NewAPI(clientID string) *API {
	return &API{
		clientID: clientID,
	}
}

// GetUser returns user for token
func (a *API) GetUser(accessToken string) (*User, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/user", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", a.clientID)
	req.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := &User{}

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
