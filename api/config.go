package api

import (
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
)

// ConfigAggregate struct holds all necessary app configuration parameters.
type ConfigAggregate struct {
	Config Config
	Rate   rate.ProviderConfig
	Sender subscription.SenderConfig
	Repo   subscription.RepoConfig
}

type Config struct {
	Name    string
	Path    string
	Version string
	Origin  string
}
