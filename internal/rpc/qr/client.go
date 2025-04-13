package qr

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

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

func (s *Service) CreateQR(ctx context.Context) (*qrproto.CreateQROut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))

	resp, err := s.client.CreateQR(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to create QR in rpc: %v", err)
	}

	return resp, nil
}

func (s *Service) VerifyAccess(ctx context.Context, token, action string) (*qrproto.VerifyQROut, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))

	req := qrproto.VerifyQRIn{
		Token:  token,
		Action: action,
	}
	resp, err := s.client.VerifyQR(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify QR in rpc: %v", err)
	}

	return resp, nil
}
