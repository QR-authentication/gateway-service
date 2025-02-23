package auth

import (
	"context"
	"fmt"
	"log"

	auth "github.com/QR-authentication/auth-proto/auth-proto"
	"github.com/QR-authentication/gateway-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Service struct {
	client auth.AuthServiceClient
}

func NewService(cfg *config.Config) *Service {
	connStr := fmt.Sprintf("%s:%s", cfg.Auth.Host, cfg.Auth.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc connection: %v", err)
	}

	client := auth.NewAuthServiceClient(conn)

	return &Service{client: client}
}

func (s *Service) Login(ctx context.Context, login, password string) (*auth.LoginOut, error) {
	resp, err := s.client.Login(ctx, &auth.LoginIn{
		Login:    login,
		Password: password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, status.Error(codes.InvalidArgument, "failed to invalid login or password")
			case codes.Internal:
				return nil, status.Errorf(codes.Internal, "failed to internal server error: %v", st.Message())
			default:
				return nil, status.Errorf(codes.Unknown, "failed to unexpected error: %v", st.Message())
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to connect to auth service: %v", err)
	}

	return resp, nil
}
