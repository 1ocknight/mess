package grpc

import (
	"context"
	"fmt"
	"tokenissuer/internal/service"
	pb "tokenissuer/internal/transport/grpc/pb/tokenissuerpb"

	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedTokenVerifierServer
	srv service.Verify
}

func NewService(srv service.Verify) *Service {
	return &Service{
		srv: srv,
	}
}

func Register(gRPCServer *grpc.Server, srv *Service) {
	pb.RegisterTokenVerifierServer(gRPCServer, srv)
}

func (h *Service) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	if req.GetTokenType() == "" {
		return nil, fmt.Errorf("not have token type")
	}
	if req.GetAccessToken() == "" {
		return nil, fmt.Errorf("not have access token")
	}

	user, err := h.srv.VerifyToken(ctx, req.GetTokenType(), req.GetAccessToken())
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	return &pb.VerifyResponse{
		SubjectId: user.ID,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}
