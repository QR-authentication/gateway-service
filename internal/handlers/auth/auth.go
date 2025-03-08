package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/QR-authentication/gateway-service/internal/config"
)

type Handler struct {
	aS         UseCase
	SigningKey string
}

func New(cfg *config.Config, aS UseCase) *Handler {
	return &Handler{aS: aS, SigningKey: cfg.Platform.SigningKey}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Login    string `json:"login"    validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var req loginRequest

	body, err := io.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		http.Error(w, "failed to invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "failed to invalid JSON format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	jwtResp, err := h.aS.Login(ctx, req.Login, req.Password)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			http.Error(w, st.Message(), grpcStatusToHTTP(st.Code()))
		} else {
			http.Error(w, "authentication failed", http.StatusForbidden)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "FA_AUTH_TOKEN",
		Value:    jwtResp.Token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func grpcStatusToHTTP(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (h *Handler) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "FA_AUTH_TOKEN",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (h *Handler) CheckJWT(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("FA_AUTH_TOKEN")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"valid": false})
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.SigningKey), nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"valid": false})
		return
	}

	if !token.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"valid": false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}

func AttachAuthRoutes(r chi.Router, handler *Handler) {
	r.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/login", handler.Login)
		authRouter.Get("/logout", handler.Logout)
		authRouter.Get("/check-jwt", handler.CheckJWT)
	})
}
