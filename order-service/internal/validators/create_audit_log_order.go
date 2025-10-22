package validators

import (
	"fmt"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
)

func ValidateAuditLogOrderBatchCreateRequest(req *pb.AuditLogOrderBatchCreateRequest) ValidationErrors {
	errs := make(ValidationErrors)

	if len(req.Orders) == 0 {
		errs["orders"] = "at least one is required"
		return errs
	}

	for i, l := range req.Orders {
		prefix := fmt.Sprintf("orders[%d]", i)
		logErrs := validateLogOrder(l, prefix)
		errs.Merge(logErrs)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateLogOrder(l *pb.LogOrder, prefix string) ValidationErrors {
	errs := make(ValidationErrors)

	if l.OrderId <= 0 {
		errs[prefix+".order_id"] = "must be greater than 0"
	}
	if l.OrderItemId <= 0 {
		errs[prefix+".order_item_id"] = "must be greater than 0"
	}
	if l.CustomerId <= 0 {
		errs[prefix+".customer_id"] = "must be greater than 0"
	}
	if l.OrderStatus == "" {
		errs[prefix+".order_status"] = "required"
	}

	return errs
}
