package api

import (
	"net/http"

	qrproto "github.com/QR-authentication/qr-proto/qr-proto"
)

type QRService interface {
	GenerateQRCode(r *http.Request) (*qrproto.CreateQROut, error)
	VerifyAccess(r *http.Request) (*qrproto.VerifyQROut, error)
}
