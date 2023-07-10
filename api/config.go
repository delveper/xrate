package api

import (
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subs"
)

// ConfigAggregate struct holds all necessary app configuration parameters.
type ConfigAggregate struct {
	Api          Config
	Rate         rate.Config
	Subscription subs.Config
}

type Config struct {
	Name    string
	Path    string
	Version string
	Origin  string
}
