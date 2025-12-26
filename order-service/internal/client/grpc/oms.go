package grpcclient

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
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

func (c *OmsGrpcClient) LogOrder(ctx context.Context, req *pb.AuditLogOrderBatchCreateRequest) (*pb.AuditLogOrderBatchCreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.AuditLogOrderBatchCreate(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *OmsGrpcClient) QueryOrders(ctx context.Context, req *pb.QueryOrdersRequest) (*pb.QueryOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.QueryOrders(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}
