package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/QR-authentication/gateway-service/internal/config"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	qS QRService
}

func New(qS QRService) *Handler {
	return &Handler{qS: qS}
}

func (h *Handler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	resp, err := h.qS.GenerateQRCode(r)
	if err != nil {
		log.Printf("failed to generate qr code: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsn, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed marshal json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsn)
}

func (h *Handler) VerifyAccess(w http.ResponseWriter, r *http.Request) {
	resp, err := h.qS.VerifyAccess(r)
	if err != nil {
		log.Printf("failed to check access ability: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsn, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed marshal json: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsn)
}

func AttachApiRoutes(r chi.Router, handler *Handler, cfg *config.Config) {
	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(func(next http.Handler) http.Handler {
			return CheckJWT(next, cfg)
		})

		apiRouter.Get("/qr", handler.GenerateQRCode)
		apiRouter.Post("/qr", handler.VerifyAccess)
	})
}
