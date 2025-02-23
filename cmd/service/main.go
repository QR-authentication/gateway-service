package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/QR-authentication/gateway-service/internal/config"
	"github.com/QR-authentication/gateway-service/internal/handlers/api"
	authhandler "github.com/QR-authentication/gateway-service/internal/handlers/auth"
	"github.com/QR-authentication/gateway-service/internal/middlewares"
	"github.com/QR-authentication/gateway-service/internal/rpc/auth"
	"github.com/QR-authentication/gateway-service/internal/rpc/qr"
	authusecase "github.com/QR-authentication/gateway-service/internal/useCase/auth"
	qrusecase "github.com/QR-authentication/gateway-service/internal/useCase/qr"
	metrics_lib "github.com/QR-authentication/metrics-lib"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	metrics, err := metrics_lib.New(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}
	defer metrics.Disconnect()

	authClient := auth.NewService(cfg)
	qrClient := qr.NewService(cfg)

	authUseCase := authusecase.New(authClient)
	qrUseCase := qrusecase.New(qrClient)

	authHandlers := authhandler.New(cfg, authUseCase)
	apiHandlers := api.New(qrUseCase)

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return middlewares.MetricMiddleware(next, metrics)
	})

	authhandler.AttachAuthRoutes(r, authHandlers)
	api.AttachApiRoutes(r, apiHandlers, cfg)

	log.Println("Server was started...")

	if err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Service.Port), r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
