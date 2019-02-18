package config

// Config app config
type Config struct {
	Host         string
	Port         int
	ClientID     string
	ClientSecret string
	Scopes       []string
	CookieSecret string
	RedirectURL  string
}
