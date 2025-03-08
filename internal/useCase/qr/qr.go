package qr

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	qrproto "github.com/QR-authentication/qr-proto/qr-proto"

	"github.com/QR-authentication/gateway-service/internal/model"
)

type Usecase struct {
	qC QRService
}

func New(qC QRService) *Usecase {
	return &Usecase{qC: qC}
}

func (uc *Usecase) GenerateQRCode(r *http.Request) (*qrproto.CreateQROut, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("failed to parse remote address:", err)
	}

	resp, err := uc.qC.CreateQR(r.Context(), ip)
	if err != nil {
		return nil, fmt.Errorf("failed to create qr in usecase: %w", err)
	}

	return resp, nil
}

func (uc *Usecase) VerifyAccess(r *http.Request) (*qrproto.VerifyQROut, error) {
	requestData := model.RequestData{}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}
	defer r.Body.Close()

	resp, err := uc.qC.VerifyAccess(r.Context(), requestData.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to check access in usecase: %w", err)
	}

	return resp, nil
}
