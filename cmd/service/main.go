package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/QR-authentication/gateway-service/internal/config"
	"github.com/QR-authentication/gateway-service/internal/handler/api"
	"github.com/QR-authentication/gateway-service/internal/rpc/qr"
	qrusecase "github.com/QR-authentication/gateway-service/internal/useCase/qr"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	//metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "gateway", cfg.Platform.Env)
	//if err != nil {
	//	log.Fatalf("failed to init metrics: %v", err)
	//}
	//defer metrics.Disconnect()

	qrClient := qr.NewService(cfg)

	qrUseCase := qrusecase.New(qrClient)

	apiHandlers := api.New(qrUseCase)

	r := chi.NewRouter()

	api.AttachApiRoutes(r, apiHandlers)

	log.Println("Server was started...")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Service.Port), r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
