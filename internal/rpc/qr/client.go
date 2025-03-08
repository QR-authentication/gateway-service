package qr

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	qrproto "github.com/QR-authentication/qr-proto/qr-proto"

	"github.com/QR-authentication/gateway-service/internal/config"
)

type Service struct {
	client qrproto.QRServiceClient
}

func NewService(cfg *config.Config) *Service {
	connStr := fmt.Sprintf("%s:%s", cfg.QR.Host, cfg.QR.Port)

	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}

	client := qrproto.NewQRServiceClient(conn)

	return &Service{client: client}
}

func (s *Service) CreateQR(ctx context.Context, ip string) (*qrproto.CreateQROut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))

	req := qrproto.CreateQRIn{
		Ip: ip,
	}

	resp, err := s.client.CreateQR(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR in rpc: %v", err)
	}

	return resp, nil
}

func (s *Service) VerifyAccess(ctx context.Context, token string) (*qrproto.VerifyQROut, error) {
	req := qrproto.VerifyQRIn{
		Token: token,
	}

	resp, err := s.client.VerifyQR(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify QR in rpc: %v", err)
	}

	return resp, nil
}
