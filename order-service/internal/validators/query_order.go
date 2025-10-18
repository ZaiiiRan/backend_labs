package validators

import (
	"fmt"

	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
)

func ValidateV1QueryOrdersRequest(req dto.V1QueryOrdersRequest) ValidationErrors {
	errs := make(map[string]string)

	if req.Page < 1 {
		errs["page"] = "must be greater than or equal to 1"
	}
	if req.PageSize < 1 {
		errs["page_size"] = "must be greater than or equal to 1"
	}
	if req.PageSize > 100 {
		errs["page_size"] = "must be less than or equal to 100"
	}

	for i, o := range req.IDs {
		key := fmt.Sprintf("ids[%d]", i)
		if o <= 0 {
			errs[key] = "must be greater than 0"
		}
	}

	for i, cId := range req.CustomerIDs {
		key := fmt.Sprintf("customer_ids[%d]", i)
		if cId <= 0 {
			errs[key] = "must be greater than 0"
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
