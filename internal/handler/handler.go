package handler

import (
	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"account-service/docs"
	"account-service/internal/config"
	"account-service/internal/handler/http"
	"account-service/internal/service/auth"
	"account-service/pkg/server/router"
)

type Dependencies struct {
	AuthService *auth.Service

	Configs config.Configs
}

// Configuration is an alias for a function that will take in a pointer to a Handler and modify it
type Configuration func(h *Handler) error

// Handler is an implementation of the Handler
type Handler struct {
	dependencies Dependencies

	HTTP *chi.Mux
}

// New takes a variable amount of Configuration functions and returns a new Handler
// Each Configuration will be called in the order they are passed in
func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	// Insert the handler
	h = &Handler{
		dependencies: d,
	}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()

		reg := prometheus.NewRegistry()

		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)

		h.HTTP.Use(middleware.Timeout(h.dependencies.Configs.APP.Timeout))

		prometheusMiddleware := chiprometheus.NewMiddleware("product-service")

		h.HTTP.Use(prometheusMiddleware)

		// Init swagger handler
		docs.SwaggerInfo.BasePath = h.dependencies.Configs.APP.Path
		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

		// Init service handlers
		authHandler := http.NewAuthHandler(h.dependencies.AuthService)
		
		h.HTTP.Route("/", func(r chi.Router) {
			r.Mount("/auth", authHandler.Routes())
		})

		return
	}
}
