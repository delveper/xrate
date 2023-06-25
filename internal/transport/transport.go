// Package transport provides functionality to handle HTTP requests.
package transport

import (
	"net/http"

	"github.com/GenesisEducationKyiv/main-project-delveper/internal/rate"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subscription"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/filestore"
	"github.com/GenesisEducationKyiv/main-project-delveper/sys/logger"
)

// API struct contains configuration details and logger instance.
type API struct {
	cfg Config
	log *logger.Logger
}

// Config struct holds all necessary configuration parameters.
type Config struct {
	DBPath       string
	EmailAPIkey  string
	EmailAddress string
	RateEndpoint string
}

// New returns a new API instance with provided configuration and logger.
func New(cfg Config, log *logger.Logger) *API {
	return &API{
		cfg: cfg,
		log: log,
	}
}

// Handle return http.Handler with all application routes defined.
func (a *API) Handle() http.Handler {
	rateSvc := rate.NewService(a.cfg.RateEndpoint)
	rateHdl := rate.NewHandler(rateSvc, a.log)

	emailStore := filestore.New[subscription.Email](a.cfg.DBPath)
	emailRepo := subscription.NewRepo(emailStore)
	senderSvc := subscription.NewSender(a.cfg.EmailAddress, a.cfg.EmailAPIkey)
	subscriptionSvc := subscription.NewService(emailRepo, rateSvc, senderSvc)
	subscriptionHdl := subscription.NewHandler(subscriptionSvc, a.log)

	mux := http.NewServeMux()

	// TODO: In further iterations add 3d party router.
	mux.Handle("/api/rate", a.WithMethod(http.MethodGet)(rateHdl.Rate))
	mux.Handle("/api/subscribe", a.WithMethod(http.MethodPost)(subscriptionHdl.Subscribe))
	mux.Handle("/api/sendEmails", a.WithMethod(http.MethodPost)(subscriptionHdl.SendEmails))

	hdl := ChainMiddlewares(mux.ServeHTTP,
		a.WithLogRequest,
		a.WithJSON,
		a.WithoutPanic,
	)

	return hdl
}
