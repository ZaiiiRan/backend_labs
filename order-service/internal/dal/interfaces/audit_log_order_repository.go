package interfaces

import (
	"context"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
)

type AuditLogOrderRepository interface {
	BulkInsert(ctx context.Context, items []models.V1AuditLogOrderDal) ([]models.V1AuditLogOrderDal, error)
}
