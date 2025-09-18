package interfaces

import (
	"context"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
)

type OrderRepository interface {
	BulkInsert(ctx context.Context, orders []models.V1OrderDal) ([]models.V1OrderDal, error)
	Query(ctx context.Context, query models.QueryOrderItemsDalModel) ([]models.V1OrderDal, error)
}
