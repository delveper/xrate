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
	RapidApi, CoinApi, Ninjas, AlphaVantage, CoinYep ProviderConfig
	ClientRetryMax                                   int
}

type ProviderConfig struct {
	Endpoint string
	Key      string
}

type EmailConfig struct {
	SenderAddress string
	SenderKey     string
}

type SubscriptionConfig struct {
	Data string
}
