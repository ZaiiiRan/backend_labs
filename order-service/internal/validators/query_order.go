package validators

import (
	"fmt"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
)

func ValidateQueryOrdersRequest(req *pb.QueryOrdersRequest) ValidationErrors {
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

	for i, o := range req.Ids {
		key := fmt.Sprintf("ids[%d]", i)
		if o <= 0 {
			errs[key] = "must be greater than 0"
		}
	}

	for i, cId := range req.CustomerIds {
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
