package api

// Config struct holds all necessary app configuration parameters.
type Config struct {
	ApiConfig          ApiConfig
	RateConfig         RateConfig
	EmailConfig        EmailConfig
	SubscriptionConfig SubscriptionConfig
}

type ApiConfig struct {
	Name    string
	Path    string
	Version string
	Origin  string
}

type RateConfig struct {
	Endpoint string
	RetryMax int
}

type EmailConfig struct {
	SenderAddress string
	SenderKey     string
}

type SubscriptionConfig struct {
	Data string
}
