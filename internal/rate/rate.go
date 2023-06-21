/*
Package rate provides a client for fetching exchange BTC2UAH rates from a specified endpoint.

The `Service` struct represents the rate service and includes the following fields:
- `Endpoint`: The URL endpoint for fetching exchange rates.
- `Client`: The HTTP client used for making requests.

The package includes the following main functions:
- `NewService(endpoint string) *Service`: Creates a new rate service instance with the specified endpoint URL.
- `Get() (float64, error)`: Retrieves the exchange rate from the endpoint and returns it as a float64 value.
*/
package rate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	Endpoint string
	Client   *http.Client
}

func NewService(endpoint string) *Service {
	return &Service{
		Endpoint: endpoint,
		Client:   new(http.Client),
	}
}

func (a *Service) Get() (float64, error) {
	req, err := http.NewRequest(http.MethodGet, a.Endpoint, nil)
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

	var rate struct {
		Rates struct{ UAH struct{ Value float64 } }
	}

	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	return rate.Rates.UAH.Value, nil
}
