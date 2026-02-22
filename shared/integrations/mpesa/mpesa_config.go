package mpesa

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	Passkey        string
	Shortcode      string
	CallbackURL    string
	Environment    string // sandbox | production
}

func (c *Config) BaseURL() string {
	if c.Environment == "production" {
		return "https://api.safaricom.co.ke"
	}
	return "https://sandbox.safaricom.co.ke"
}
