package rate

type Config struct {
	Provider struct {
		ExchangeRateHost, CoinYep, RapidApi, Ninjas, AlphaVantage, CoinApi ProviderConfig
	}
	Client struct {
		RetryMax int
	}
}

type ProviderConfig struct {
	Endpoint string
	Header   string
	Key      string
}
