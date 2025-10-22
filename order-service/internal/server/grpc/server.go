package grpcserver

import (
	"context"
	"fmt"
	"net"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/server/grpc/services"
	"google.golang.org/grpc"
)

type Server struct {
	srv          *grpc.Server
	listener     net.Listener
	orderService *services.OrderService
}

func NewServer(port int, orderService *services.OrderService) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, orderService)

	return &Server{
		srv:          s,
		listener:     lis,
		orderService: orderService,
	}, nil
}

func (s *Server) Start() error {
	return s.srv.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}
