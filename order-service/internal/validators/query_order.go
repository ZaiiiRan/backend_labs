package validators

import "github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"

func ValidateV1QueryOrdersRequest(req dto.V1QueryOrdersRequest) ValidationErrors {
	errs := make(map[string]string)

	if req.Page < 1 {
		errs["page"] = "must be greater than or equal to 1"
	}
	if req.PageSize < 1 {
		errs["pageSize"] = "must be greater than or equal to 1"
	}
	if req.PageSize > 100 {
		errs["pageSize"] = "must be less than or equal to 100"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
