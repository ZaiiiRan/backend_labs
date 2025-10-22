package mappers

import (
	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
)

func DtoQueryOrderItemsToBll(q dto.V1QueryOrdersRequest) bll.QueryOrderItemsModel {
	return bll.QueryOrderItemsModel{
		IDs:               q.IDs,
		CustomerIDs:       q.CustomerIDs,
		Page:              q.Page,
		PageSize:          q.PageSize,
		IncludeOrderItems: q.IncludeOrderItems,
	}
}

func PbQueryOrderItemsToBll(q *pb.QueryOrdersRequest) bll.QueryOrderItemsModel {
	return bll.QueryOrderItemsModel{
		IDs:               q.Ids,
		CustomerIDs:       q.CustomerIds,
		Page:              int(q.Page),
		PageSize:          int(q.PageSize),
		IncludeOrderItems: q.IncludeOrderItems,
	}
}
