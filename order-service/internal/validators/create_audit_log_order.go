package validators

import (
	"fmt"

	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
)

func ValidateV1CreateAuditLogOrderRequest(req dto.V1CreateAuditLogOrderRequest) ValidationErrors {
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

func validateLogOrder(l dto.V1LogOrder, prefix string) ValidationErrors {
	errs := make(ValidationErrors)

	if l.OrderId <= 0 {
		errs[prefix+".orderId"] = "must be greater than 0"
	}
	if l.OrderItemId <= 0 {
		errs[prefix+".orderItemId"] = "must be greater than 0"
	}
	if l.CustomerId <= 0 {
		errs[prefix+".customerId"] = "must be greater than 0"
	}
	if l.OrderStatus == "" {
		errs[prefix+".orderStatus"] = "required"
	}

	return errs
}
