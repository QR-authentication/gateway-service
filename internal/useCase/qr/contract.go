package qr

import (
	"context"

	qrproto "github.com/QR-authentication/qr-proto/qr-proto"
)

type QRService interface {
	CreateQR(ctx context.Context, uuid, ip string) (*qrproto.CreateQROut, error)
	VerifyAccess(ctx context.Context, token string) (*qrproto.VerifyQROut, error)
}
