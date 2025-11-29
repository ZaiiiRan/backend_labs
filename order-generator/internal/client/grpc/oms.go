package grpcclient

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	config "github.com/ZaiiiRan/backend_labs/order-generator/internal/config/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OmsGrpcClient struct {
	conn   *grpc.ClientConn
	client pb.OrderServiceClient
}

func NewOmsGrpcClient(cfg config.GrpcClientSettings) (*OmsGrpcClient, error) {
	conn, err := grpc.Dial(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %w", err)
	}

	client := pb.NewOrderServiceClient(conn)
	return &OmsGrpcClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *OmsGrpcClient) Close() error {
	return c.conn.Close()
}

func (c *OmsGrpcClient) BatchCreate(ctx context.Context, req *pb.BatchCreateRequest) (*pb.BatchCreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.BatchCreate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc call BatchCreate failed: %w", err)
	}

	return resp, nil
}
