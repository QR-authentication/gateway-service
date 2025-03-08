package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/QR-authentication/gateway-service/internal/config"
	"github.com/QR-authentication/gateway-service/internal/model"
)

func CheckJWT(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("checking JWT for: %s", r.URL.Path)
		cookie, err := r.Cookie("FA_AUTH_TOKEN")
		if err != nil {
			log.Printf("failed to get cookie: %v", err)
			invalidateCookie(w, r.TLS != nil)
			http.Error(w, "failed to missing or invalid cookie", http.StatusUnauthorized)
			return
		}

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("failed to unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.Platform.SigningKey), nil
		})

		if err != nil {
			log.Printf("failed to parse JWT: %v", err)
			invalidateCookie(w, r.TLS != nil)
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "failed to token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "failed to invalid token", http.StatusUnauthorized)
			}
			return
		}

		if !token.Valid {
			log.Printf("invalid JWT token")
			invalidateCookie(w, r.TLS != nil)
			http.Error(w, "failed to token is invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), config.KeyUUID, claims.UUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func invalidateCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "FA_AUTH_TOKEN",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}
