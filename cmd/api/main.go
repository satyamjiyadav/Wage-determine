package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "prevailing-wage-service/internal/adapter/http"
	"prevailing-wage-service/internal/infrastructure/cache"
	"prevailing-wage-service/internal/infrastructure/persistence"
	"prevailing-wage-service/internal/usecase"
)

func main() {
	log.Println("==========================================================")
	log.Println(" Starting U.S. DOL Prevailing Wage Service (Golang)      ")
	log.Println("==========================================================")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Initialize Persistence Infrastructure & Cache
	memoryRepo := persistence.NewMemoryRepository()
	memoryCache := cache.NewMemoryCache()

	// 2. Initialize Application Use Cases
	wageLookupUC := usecase.NewWageLookupUseCase(memoryRepo, memoryRepo, memoryRepo, memoryCache)
	levelCalcUC := usecase.NewLevelCalculatorUseCase(memoryRepo, memoryRepo, memoryRepo, memoryRepo, memoryRepo)
	discoveryUC := usecase.NewDiscoveryUseCase(memoryRepo, memoryRepo)
	auditUC := usecase.NewDeterminationAuditUseCase(memoryRepo)

	// 3. Initialize HTTP Handlers
	wageHandler := httpAdapter.NewWageHandler(wageLookupUC, levelCalcUC)
	discoveryHandler := httpAdapter.NewDiscoveryHandler(discoveryUC)
	auditHandler := httpAdapter.NewAuditHandler(auditUC)
	healthHandler := httpAdapter.NewHealthHandler()

	// 4. Initialize HTTP Router
	router := httpAdapter.NewRouter(wageHandler, discoveryHandler, auditHandler, healthHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start Server in Goroutine
	go func() {
		log.Printf("[INFO] Server listening on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] HTTP server failed to start: %v\n", err)
		}
	}()

	// 6. Graceful Shutdown Signal Handler
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[INFO] Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v\n", err)
	}

	fmt.Println("[SUCCESS] Service stopped gracefully.")
}
