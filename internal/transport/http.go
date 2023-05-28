package transport

import (
	"github.com/delveper/gentest/internal/rate"
	"github.com/delveper/gentest/internal/subscription"
	"github.com/delveper/gentest/sys/filestore"
	"github.com/delveper/gentest/sys/logger"
	"net/http"
)

type API struct {
	cfg Config
	log *logger.Logger
}

type Config struct {
	DBPath       string
	EmailAPIkey  string
	EmailAddress string
}

func New(cfg Config, log *logger.Logger) *API {
	return &API{
		cfg: cfg,
		log: log,
	}
}

// Handle return http.Handler with all application routes defined.
func (a *API) Handle() http.Handler {
	rateSvc := rate.NewService()
	rateHdl := rate.NewHandler(rateSvc, a.log)

	emailStore := filestore.New[subscription.Email](a.cfg.DBPath)
	emailRepo := subscription.NewRepo(emailStore)
	senderSvc := subscription.NewSender(a.cfg.EmailAPIkey, a.cfg.EmailAddress)
	subscriptionSvc := subscription.NewService(emailRepo, rateSvc, senderSvc)
	subscriptionHdl := subscription.NewHandler(subscriptionSvc, a.log)

	mux := http.NewServeMux()
	mux.Handle("/api/rate", a.WithMethod(http.MethodGet)(rateHdl.Rate))
	mux.Handle("/api/subscribe", a.WithMethod(http.MethodPost)(subscriptionHdl.Subscribe))
	mux.Handle("/api/sendEmails", a.WithMethod(http.MethodPost)(subscriptionHdl.SendEmails))

	hdl := ChainMiddlewares(mux.ServeHTTP,
		a.WithLogRequest,
		a.WithCORS,
		a.WithJSON,
		a.WithoutPanic,
	)

	return hdl
}
