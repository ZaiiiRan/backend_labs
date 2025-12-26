package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	grpcclient "github.com/ZaiiiRan/backend_labs/order-service/internal/client/grpc"
	dalconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/utils"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type OrderCreatedMessageProcessor struct {
	client *grpcclient.OmsGrpcClient
	log    *zap.SugaredLogger
}

func NewOrderCreatedMessageProcessor(client *grpcclient.OmsGrpcClient, log *zap.SugaredLogger) dalconsumer.MessageProcessor {
	return &OrderCreatedMessageProcessor{
		client: client,
		log:    log,
	}
}

func (p *OrderCreatedMessageProcessor) ProcessMessage(ctx context.Context, batch []dalconsumer.MessageInfo) (bool, error) {
	var orders []messages.OrderCreatedMessage
	for _, msg := range batch {
		var o messages.OrderCreatedMessage
		if err := json.Unmarshal(msg.Body, &o); err != nil {
			p.log.Errorw("order_created_message_processor.unmarshal_failed", "err", err, "body", string(msg.Body))
			return false, fmt.Errorf("unmarshal: %w", err)
		}
		orders = append(orders, o)
	}

	req := &pb.AuditLogOrderBatchCreateRequest{}
	for _, order := range orders {
		for _, item := range order.OrderItems {
			log := &pb.LogOrder{
				OrderId:     order.Id,
				OrderItemId: item.Id,
				CustomerId:  order.CustomerID,
				OrderStatus: models.ORDER_STATUS_CREATED.String(),
			}
			p.log.Infow("order_created_message_processor.log_order", "log_order", log)
			req.Orders = append(req.Orders, log)
		}
	}

	if _, err := p.client.LogOrder(ctx, req); err != nil {
		p.log.Errorw("order_created_message_processor.grpc_call_failed", "err", err)

		needToRequeue := false
		st, err := utils.GetGrpcErrStatus(err)
		if err != nil || st.Code() != codes.InvalidArgument {
			needToRequeue = true
		}

		return needToRequeue, fmt.Errorf("grpc: %w", err)
	}

	p.log.Infow("order_created_message_processor.batch_processed", "count", len(orders))
	return false, nil
}
