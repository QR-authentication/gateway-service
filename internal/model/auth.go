package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UUID      string
	ExpiresAt time.Time
	IssuedAt  time.Time
	jwt.RegisteredClaims
}
