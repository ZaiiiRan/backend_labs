package services

import (
	"context"
	"time"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/mappers"
	bll "github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"
	"github.com/ZaiiiRan/backend_labs/order-service/internal/dal/interfaces"
	dal "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/models"
	unitofwork "github.com/ZaiiiRan/backend_labs/order-service/internal/dal/unit_of_work/postgres"
	"go.uber.org/zap"
)

type AuditLogOrderService struct {
	uow                   *unitofwork.UnitOfWork
	auditLogOrderItemRepo interfaces.AuditLogOrderRepository
	log                   *zap.SugaredLogger
}

func NewAuditLogOrderService(
	uow *unitofwork.UnitOfWork,
	auditLogOrderItemRepo interfaces.AuditLogOrderRepository,
	log *zap.SugaredLogger,
) *AuditLogOrderService {
	return &AuditLogOrderService{
		uow:                   uow,
		auditLogOrderItemRepo: auditLogOrderItemRepo,
		log:                   log,
	}
}

func (s *AuditLogOrderService) BatchInsert(ctx context.Context, logs []bll.AuditLogOrder) ([]bll.AuditLogOrder, error) {
	now := time.Now().UTC()
	s.log.Infow("audit_log_order_service.batch_insert_start", "logs_count", len(logs))

	_, err := s.uow.BeginTransaction(ctx)
	if err != nil {
		s.log.Errorw("audit_log_order_service.begin_transaction_failed", "err", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			s.uow.Rollback(ctx)
			s.log.Warnw("audit_log_order_service.transaction_rollback", "err", err)
		}
	}()

	var dalLogs []dal.V1AuditLogOrderDal
	for _, l := range logs {
		d := mappers.BllAuditLogOrderToDal(l)
		d.CreatedAt = now
		d.UpdatedAt = now
		dalLogs = append(dalLogs, d)
	}

	insertedLogs, err := s.auditLogOrderItemRepo.BulkInsert(ctx, dalLogs)
	if err != nil {
		s.log.Errorw("audit_log_order_service.bulk_insert_logs_failed", "err", err)
		return nil, err
	}

	if err := s.uow.Commit(ctx); err != nil {
		s.log.Errorw("audit_log_order_service.commit_transaction_failed", "err", err)
		return nil, err
	}

	var result []bll.AuditLogOrder
	for _, l := range insertedLogs {
		result = append(result, mappers.DalAuditLogOrderToBll(l))
	}

	s.log.Infow("audit_log_order_service.batch_insert_success", "inserted_logs_count", len(result))
	return result, nil
}

func (s *AuditLogOrderService) UnitOfWork() *unitofwork.UnitOfWork {
	return s.uow
}
