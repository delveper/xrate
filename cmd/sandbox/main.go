package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/GenesisEducationKyiv/main-project-delveper/sys/env"
)

type Config struct {
	U url.URL `default:"https://api.coingecko.com/api/v3/"`
}

func main() {
	var cfg Config
	if err := env.ParseTo(".env", &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.U)
}
