package grpcgateway

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	srv *http.Server
}

func NewServer(ctx context.Context, port int, grpcPort int) (*Server, error) {
	mux := runtime.NewServeMux()
	grpcAddr := fmt.Sprintf("localhost:%d", grpcPort)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register gateway handler: %w", err)
	}

	swaggerDir := filepath.Join("gen", "openapiv2", "order-service", "v1")

	rootMux := http.NewServeMux()
	rootMux.Handle("/", mux)

	rootMux.Handle("/swagger/", http.StripPrefix("/swagger/",
		http.FileServer(http.Dir(swaggerDir)),
	))

	rootMux.Handle("/docs/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/order_service.swagger.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
		),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rootMux,
	}

	return &Server{srv: srv}, nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.srv.Addr
}
