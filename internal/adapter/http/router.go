package http

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	wageHandler *WageHandler,
	discoveryHandler *DiscoveryHandler,
	auditHandler *AuditHandler,
	healthHandler *HealthHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// CORS Setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Root Welcome & Interactive API Documentation Landing Page
	r.Get("/", ServeLandingPage)

	// System & Health Endpoints
	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/metrics", healthHandler.Metrics)

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Wage Lookup & Level Assessment Endpoints
		r.Get("/wages/lookup", wageHandler.LookupWage)
		r.Post("/wages/determine-level", wageHandler.DetermineWageLevel)
		r.Post("/wages/batch-lookup", wageHandler.BatchLookupWages)

		// Occupation & Location Discovery Endpoints
		r.Get("/occupations/search", discoveryHandler.SearchOccupations)
		r.Get("/occupations/{soc_code}", discoveryHandler.GetOccupationDetails)
		r.Get("/locations/resolve", discoveryHandler.ResolveLocation)

		// Compliance Audit Trail Endpoints
		r.Get("/determinations/{determination_number}", auditHandler.GetDetermination)
	})

	log.Println("[INFO] Registered prevailing wage REST API routes successfully")
	return r
}
