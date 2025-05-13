package web

import (
	"log"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	LegalEntitiesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "legal_entities_requests_total",
		Help: "Total number of requests to /legal-entities",
	})
	BankAccountsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bank_accounts_requests_total",
		Help: "Total number of requests to /bank-accounts",
	})
)

func initMetricsRoutes(a *Web, e *echo.Echo) {
	// Регистрация через DefaultRegisterer
	prometheus.MustRegister(LegalEntitiesCounter)
	prometheus.MustRegister(BankAccountsCounter)

	e.Use(echoprometheus.NewMiddleware(a.Options.APP_NAME))

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			path := c.Path()
			log.Println("➡️ full URI:", c.Request().URL.Path, "| route:", path)

			if strings.Contains(path, "/legal-entities") {
				log.Println("🟢 /legal-entities matched – counter++")
				LegalEntitiesCounter.Inc()
			}
			if strings.Contains(path, "/bank-accounts") {
				log.Println("🟢 /bank-accounts matched – counter++")
				BankAccountsCounter.Inc()
			}
		},
	}))

	e.GET("/metrics", echoprometheus.NewHandler())
}
