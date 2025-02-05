package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/QR-authentication/gateway-service/internal/config"
	"github.com/QR-authentication/gateway-service/internal/handler/api"
	"github.com/QR-authentication/gateway-service/internal/middlewares"
	"github.com/QR-authentication/gateway-service/internal/rpc/qr"
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

	qrClient := qr.NewService(cfg)

	qrUseCase := qrusecase.New(qrClient)

	apiHandlers := api.New(qrUseCase)

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return middlewares.MetricMiddleware(next, metrics)
	})

	api.AttachApiRoutes(r, apiHandlers)

	log.Println("Server was started...")

	if err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Service.Port), r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
