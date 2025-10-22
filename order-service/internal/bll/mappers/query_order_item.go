package mappers

import (
	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
)

func PbQueryOrderItemsToBll(q *pb.QueryOrdersRequest) bll.QueryOrderItemsModel {
	return bll.QueryOrderItemsModel{
		IDs:               q.Ids,
		CustomerIDs:       q.CustomerIds,
		Page:              int(q.Page),
		PageSize:          int(q.PageSize),
		IncludeOrderItems: q.IncludeOrderItems,
	}
}
