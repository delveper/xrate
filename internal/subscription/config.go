package subscription

type Config struct {
	Sender SenderConfig
	Repo   RepoConfig
}

type SenderConfig struct {
	Address string
	Key     string
}

type RepoConfig struct {
	Data string
}
