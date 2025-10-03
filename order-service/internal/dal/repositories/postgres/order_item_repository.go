package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/interfaces"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
)

type OrderItemRepository struct {
	uow *unitofwork.UnitOfWork
}

func NewOrderItemRepository(uow *unitofwork.UnitOfWork) interfaces.OrderItemRepository {
	return &OrderItemRepository{uow: uow}
}

func (r *OrderItemRepository) BulkInsert(ctx context.Context, items []models.V1OrderItemDal) ([]models.V1OrderItemDal, error) {
	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Conn().LoadType(ctx, "v1_order_item"); err != nil {
		return nil, fmt.Errorf("load type v1_order_item: %w", err)
	}

	sql := `
		insert into order_items (
			order_id,
			product_id,
			quantity,
			product_title,
			product_url,
			price_cents,
			price_currency,
			created_at,
			updated_at
		)
		select
			order_id,
			product_id,
			quantity,
			product_title,
			product_url,
			price_cents,
			price_currency,
			created_at,
			updated_at
		from unnest($1::v1_order_item[])
		returning
			id,
			order_id,
			product_id,
			quantity,
			product_title,
			product_url,
			price_cents,
			price_currency,
			created_at,
			updated_at;
	`

	rows, err := conn.Query(ctx, sql, items)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.V1OrderItemDal
	for rows.Next() {
		var i models.V1OrderItemDal
		if err := rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.ProductTitle,
			&i.ProductURL, &i.PriceCents, &i.PriceCurr, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}

	return result, nil
}

func (r *OrderItemRepository) Query(ctx context.Context, q models.QueryOrderItemsDalModel) ([]models.V1OrderItemDal, error) {
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
			order_id,
			product_id,
			quantity,
			product_title,
			product_url,
			price_cents,
			price_currency,
			created_at,
			updated_at
		from order_items
	`)

	if len(q.IDs) > 0 {
		where = append(where, fmt.Sprintf("id = any($%d)", argPos))
		args = append(args, q.IDs)
		argPos++
	}
	if len(q.OrderIDs) > 0 {
		where = append(where, fmt.Sprintf("order_id = any($%d)", argPos))
		args = append(args, q.OrderIDs)
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

	var result []models.V1OrderItemDal
	for rows.Next() {
		var i models.V1OrderItemDal
		if err := rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.ProductTitle,
			&i.ProductURL, &i.PriceCents, &i.PriceCurr, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}

	return result, rows.Err()
}
