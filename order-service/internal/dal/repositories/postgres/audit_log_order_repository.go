package repositories

import (
	"context"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/interfaces"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
)

type AuditLogOrderRepository struct {
	uow *unitofwork.UnitOfWork
}

func NewAuditLogOrderRepository(uow *unitofwork.UnitOfWork) interfaces.AuditLogOrderRepository {
	return &AuditLogOrderRepository{uow: uow}
}

func (r *AuditLogOrderRepository) BulkInsert(ctx context.Context, items []models.V1AuditLogOrderDal) ([]models.V1AuditLogOrderDal, error) {
	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	sql := `
		insert into audit_log_order (
			order_id,
			order_item_id,
			customer_id,
			order_status,
			created_at,
			updated_at
		)
		select
			(i).order_id,
			(i).order_item_id,
			(i).customer_id,
			(i).order_status,
			(i).created_at,
			(i).updated_at
		from unnest($1::v1_audit_log_order[]) as i
		returning
			id,
			order_id,
			order_item_id,
			customer_id,
			order_status,
			created_at,
			updated_at;
	`

	rows, err := conn.Query(ctx, sql, items)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.V1AuditLogOrderDal
	for rows.Next() {
		var i models.V1AuditLogOrderDal
		if err := rows.Scan(&i.ID, &i.OrderID, &i.OrderItemID, &i.CustomerID,
			&i.OrderStatus, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}

	return result, nil
}
