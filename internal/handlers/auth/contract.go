package auth

import (
	"context"

	auth "github.com/QR-authentication/auth-proto/auth-proto"
)

type AuthService interface {
	Login(ctx context.Context, login string, password string) (*auth.LoginOut, error)
}
