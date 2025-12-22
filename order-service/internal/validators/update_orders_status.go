package validators

import (
	"fmt"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/bll/models"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
)

func ValidateUpdateOrdersStatusRequest(req *pb.UpdateOrdersStatusRequest) ValidationErrors {
	errs := make(ValidationErrors)

	if len(req.OrderIds) == 0 {
		errs["order_ids"] = "at least one is required"
	}

	parsedStatus := models.StringToOrderStatus(req.NewStatus)
	if parsedStatus == "" {
		errs["new_status"] = "unknown status"
	}

	for i, oId := range req.OrderIds {
		if oId <= 0 {
			prefix := fmt.Sprintf("order_ids[%d]", i)
			errs[prefix] = "must be greater than 0"
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
