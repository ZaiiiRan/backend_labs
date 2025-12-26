package generator

import (
	"context"
	"math/rand"
	"time"

	grpcclient "github.com/ZaiiiRan/backend_labs/order-generator/internal/client/grpc"
	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
)

type Generator struct {
	client      *grpcclient.OmsGrpcClient
	customerIDs []int64
}

func NewGenerator(client *grpcclient.OmsGrpcClient, customerIDs []int64) *Generator {
	return &Generator{
		client:      client,
		customerIDs: customerIDs,
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
			CustomerId:         g.randomCustomerID(),
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
		g.UpdateOrdersStatusRandom(ctx, resp.Orders)
	}
}

func (g *Generator) UpdateOrdersStatusRandom(ctx context.Context, orders []*pb.Order) error {
	if len(orders) == 0 {
		return nil
	}

	k := rand.Intn(len(orders)) + 1

	shuffled := make([]*pb.Order, len(orders))
	copy(shuffled, orders)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	ids := make([]int64, 0, k)
	for i := range k {
		ids = append(ids, shuffled[i].Id)
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

func (g *Generator) randomCustomerID() int64 {
	return g.customerIDs[rand.Intn(len(g.customerIDs))]
}
