package rate

type Config struct {
	Provider struct {
		ExchangeRateHost, Ninjas, AlphaVantage, CoinApi, CoinYep ProviderConfig
	}
	Client struct {
		RetryMax int
	}
}

type ProviderConfig struct {
	Name     string
	Endpoint string
	Header   string
	Key      string
}
