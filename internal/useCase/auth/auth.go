package auth

import (
	"context"

	auth "github.com/QR-authentication/auth-proto/auth-proto"
)

type Usecase struct {
	aC AuthClient
}

func New(aC AuthClient) *Usecase {
	return &Usecase{aC: aC}
}

func (uc *Usecase) Login(ctx context.Context, login string, password string) (*auth.LoginOut, error) {
	return uc.aC.Login(ctx, login, password)
}
