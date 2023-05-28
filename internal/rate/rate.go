package rate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct{ *http.Client }

func NewService() *Service {
	return &Service{new(http.Client)}
}

func (a *Service) Get() (float64, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.coingecko.com/api/v3/exchange_rates", nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	rate := struct {
		Rates struct {
			UAH struct {
				Value float64
			}
		}
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	return rate.Rates.UAH.Value, nil
}
