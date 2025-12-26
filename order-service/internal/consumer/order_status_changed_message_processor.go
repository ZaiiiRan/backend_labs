package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	grpcclient "github.com/ZaiiiRan/backend_labs/order-service/internal/client/grpc"
	dalconsumer "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/consumer"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/utils"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/messages"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type OrderStatusChangedMessageProcessor struct {
	client *grpcclient.OmsGrpcClient
	log    *zap.SugaredLogger
}

func NewOrderStatusChangedMessageProcessor(client *grpcclient.OmsGrpcClient, log *zap.SugaredLogger) dalconsumer.MessageProcessor {
	return &OrderStatusChangedMessageProcessor{
		client: client,
		log:    log,
	}
}

func (p *OrderStatusChangedMessageProcessor) ProcessMessage(ctx context.Context, batch []dalconsumer.MessageInfo) (bool, error) {
	orders := make([]messages.OrderStatusChangedMessage, 0, len(batch))
	ids := make([]int64, 0, len(batch))
	for _, msg := range batch {
		var o messages.OrderStatusChangedMessage
		if err := json.Unmarshal(msg.Body, &o); err != nil {
			p.log.Errorw("order_status_changed_message_processor.unmarshal_failed", "err", err, "body", string(msg.Body))
			return false, fmt.Errorf("unmarshal: %w", err)
		}
		orders = append(orders, o)
		ids = append(ids, o.OrderId)
	}

	ordersResp, err := p.client.QueryOrders(ctx, &pb.QueryOrdersRequest{
		Ids:               ids,
		IncludeOrderItems: true,
		Page:              1,
		PageSize:          int32(len(ids)),
	})
	if err != nil {
		p.log.Errorw("order_status_changed_message_processor.grpc_call_failed", "err", err)

		needToRequeue := false
		st, err := utils.GetGrpcErrStatus(err)
		if err != nil || st.Code() != codes.InvalidArgument {
			needToRequeue = true
		}

		return needToRequeue, fmt.Errorf("grpc: %w", err)
	}

	req := &pb.AuditLogOrderBatchCreateRequest{}
	for _, order := range ordersResp.Orders {
		for _, item := range order.OrderItems {
			log := &pb.LogOrder{
				OrderId:     order.Id,
				OrderItemId: item.Id,
				CustomerId:  order.CustomerId,
				OrderStatus: order.Status,
			}
			p.log.Infow("order_status_changed_message_processor.log_order", "log_order", log)
			req.Orders = append(req.Orders, log)
		}
	}

	if _, err := p.client.LogOrder(ctx, req); err != nil {
		p.log.Errorw("order_status_changed_message_processor.grpc_call_failed", "err", err)

		needToRequeue := false
		st, err := utils.GetGrpcErrStatus(err)
		if err != nil || st.Code() != codes.InvalidArgument {
			needToRequeue = true
		}

		return needToRequeue, fmt.Errorf("grpc: %w", err)
	}

	p.log.Infow("order_status_changed_message_processor.batch_processed", "count", len(orders))
	return false, nil
}
