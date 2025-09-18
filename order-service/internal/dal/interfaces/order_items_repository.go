package interfaces

import (
	"context"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
)

type OrderItemRepository interface {
	BulkInsert(ctx context.Context, items []models.V1OrderItemDal) ([]models.V1OrderItemDal, error)
	Query(ctx context.Context, query models.QueryOrderItemsDalModel) ([]models.V1OrderItemDal, error)
}
