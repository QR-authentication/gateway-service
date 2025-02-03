package qr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/QR-authentication/gateway-service/internal/model"
	qrproto "github.com/QR-authentication/qr-proto/qr-proto"
)

type Usecase struct {
	qC QRService
}

func New(qC QRService) *Usecase {
	return &Usecase{qC: qC}
}

func (uc *Usecase) GenerateQRCode(r *http.Request) (*qrproto.CreateQROut, error) {
	uuid := r.URL.Query().Get("uuid")

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("failed to parse remote address:", err)
	}

	resp, err := uc.qC.CreateQR(context.Background(), uuid, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to create qr in usecase: %w", err)
	}

	return resp, nil
}

func (uc *Usecase) VerifyAccess(r *http.Request) (*qrproto.VerifyQROut, error) {
	requestData := model.RequestData{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return nil, fmt.Errorf("failed to request body is empty")
	}

	if err = json.Unmarshal(body, &requestData); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	resp, err := uc.qC.VerifyAccess(context.Background(), requestData.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to check access in usecase: %w", err)
	}

	return resp, nil
}
