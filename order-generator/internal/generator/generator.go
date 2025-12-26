package generator

import (
	"context"
	"math/rand"
	"time"

	grpcclient "github.com/ZaiiiRan/backend_labs/order-generator/internal/client/grpc"
	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
)

type Generator struct {
	client *grpcclient.OmsGrpcClient
}

func NewGenerator(client *grpcclient.OmsGrpcClient) *Generator {
	return &Generator{
		client: client,
	}
}

func (g *Generator) Start(ctx context.Context) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.generateOrder(ctx)
		}
	}
}

func (g *Generator) generateOrder(ctx context.Context) {
	const batchSize = 100

	var orders []*pb.Order
	for i := 0; i < batchSize; i++ {
		orderItem := &pb.OrderItem{
			ProductId:     rand.Int63n(1000000),
			Quantity:      1,
			ProductTitle:  randomString(5),
			ProductUrl:    "https://example.com/item/" + randomString(8),
			PriceCents:    1000,
			PriceCurrency: "RUB",
		}

		order := &pb.Order{
			CustomerId:         rand.Int63n(1000),
			DeliveryAddress:    "г. Краснодар, ул. " + randomString(6),
			TotalPriceCents:    1000,
			TotalPriceCurrency: "RUB",
			OrderItems:         []*pb.OrderItem{orderItem},
		}

		orders = append(orders, order)
	}

	req := &pb.BatchCreateRequest{Orders: orders}

	resp, err := g.client.BatchCreate(ctx, req)
	if err == nil {
		g.UpdateOrdersStatus(ctx, resp.Orders)
	}
}

func (g *Generator) UpdateOrdersStatus(ctx context.Context, orders []*pb.Order) error {
	ids := make([]int64, 0, len(orders))
	for _, order := range orders {
		ids = append(ids, order.Id)
	}

	req := &pb.UpdateOrdersStatusRequest{
		OrderIds:  ids,
		NewStatus: "processing",
	}

	_, err := g.client.UpdateOrderStatus(ctx, req)
	return err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
