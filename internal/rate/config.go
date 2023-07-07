package rate

type Config struct {
	Provider struct {
		RapidApi, Ninjas, AlphaVantage, CoinApi, CoinYep ProviderConfig
	}
	Client struct {
		RetryMax int
	}
}

type ProviderConfig struct {
	Endpoint string
	Key      string
}
