package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/interfaces"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
)

type OrderRepository struct {
	uow *unitofwork.UnitOfWork
}

func NewOrderRepository(uow *unitofwork.UnitOfWork) interfaces.OrderRepository {
	return &OrderRepository{uow: uow}
}

func (r *OrderRepository) BulkInsert(ctx context.Context, orders []models.V1OrderDal) ([]models.V1OrderDal, error) {
	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	sql := `
		insert into orders (
			customer_id,
			delivery_address,
			total_price_cents,
			total_price_currency,
			created_at,
			updated_at,
			status
		)
		select 
			(o).customer_id,
			(o).delivery_address,
			(o).total_price_cents,
			(o).total_price_currency,
			(o).created_at,
			(o).updated_at,
			(o).status
		from unnest($1::v1_order[]) as o
		returning 
			id,
			customer_id,
			delivery_address,
			total_price_cents,
			total_price_currency,
			created_at,
			updated_at,
			status;
	`

	rows, err := conn.Conn().Query(ctx, sql, orders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.V1OrderDal
	for rows.Next() {
		var o models.V1OrderDal
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.DeliveryAddress, &o.TotalPriceCents,
			&o.TotalPriceCurr, &o.CreatedAt, &o.UpdatedAt, &o.Status,
		); err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	return result, rows.Err()
}

func (r *OrderRepository) Query(ctx context.Context, q models.QueryOrdersDalModel) ([]models.V1OrderDal, error) {
	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	var (
		sb     strings.Builder
		args   []interface{}
		where  []string
		argPos = 1
	)

	sb.WriteString(`
		select id,
			customer_id,
			delivery_address,
			total_price_cents,
			total_price_currency,
			created_at,
			updated_at,
			status
		from orders
	`)

	if len(q.IDs) > 0 {
		where = append(where, fmt.Sprintf("id = any($%d)", argPos))
		args = append(args, q.IDs)
		argPos++
	}
	if len(q.CustomerIDs) > 0 {
		where = append(where, fmt.Sprintf("customer_id = any($%d)", argPos))
		args = append(args, q.CustomerIDs)
		argPos++
	}
	if len(where) > 0 {
		sb.WriteString(" where " + strings.Join(where, " and "))
	}

	if q.Limit > 0 {
		sb.WriteString(fmt.Sprintf(" limit $%d", argPos))
		args = append(args, q.Limit)
		argPos++
	}
	if q.Offset > 0 {
		sb.WriteString(fmt.Sprintf(" offset $%d", argPos))
		args = append(args, q.Offset)
		argPos++
	}

	rows, err := conn.Conn().Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.V1OrderDal
	for rows.Next() {
		var o models.V1OrderDal
		if err := rows.Scan(
			&o.ID, &o.CustomerID, &o.DeliveryAddress, &o.TotalPriceCents,
			&o.TotalPriceCurr, &o.CreatedAt, &o.UpdatedAt, &o.Status,
		); err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	return result, rows.Err()
}
