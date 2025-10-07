package validators

import (
	"fmt"

	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
)

func ValidateV1CreateOrderRequest(req *dto.V1CreateOrderRequest) ValidationErrors {
	errs := make(ValidationErrors)

	if len(req.Orders) == 0 {
		errs["orders"] = "at least one order is required"
		return errs
	}

	for i, o := range req.Orders {
		prefix := fmt.Sprintf("orders[%d]", i)
		orderErrs := validateOrder(o, prefix)
		errs.Merge(orderErrs)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateOrder(o dto.V1Order, prefix string) ValidationErrors {
	errs := make(ValidationErrors)

	if o.CustomerID <= 0 {
		errs[prefix+".customerId"] = "must be greater than 0"
	}
	if o.DeliveryAddress == "" {
		errs[prefix+".deliveryAddress"] = "required"
	}
	if o.TotalPriceCents <= 0 {
		errs[prefix+".totalPriceCents"] = "must be greater than 0"
	}
	if o.TotalPriceCurr == "" {
		errs[prefix+".totalPriceCurrency"] = "required"
	}
	if len(o.OrderItems) == 0 {
		errs[prefix+".orderItems"] = "at least one item is required"
		return errs
	}

	var sum int64
	currencies := map[string]struct{}{}
	for j, it := range o.OrderItems {
		iprefix := fmt.Sprintf("%s.orderItems[%d]", prefix, j)
		itemErrs := validateOrderItem(it, iprefix)
		errs.Merge(itemErrs)

		sum += it.PriceCents * int64(it.Quantity)
		currencies[it.PriceCurr] = struct{}{}
	}

	if sum != o.TotalPriceCents {
		errs[prefix+".totalPriceCents"] = "must equal sum of items (priceCents * quantity)"
	}
	if len(currencies) > 1 {
		errs[prefix+".orderItems.priceCurrency"] = "all items must have the same currency"
	}
	if len(o.OrderItems) > 0 && o.OrderItems[0].PriceCurr != o.TotalPriceCurr {
		errs[prefix+".totalPriceCurrency"] = "must equal items currency"
	}

	return errs
}

func validateOrderItem(it dto.V1OrderItem, prefix string) ValidationErrors {
	errs := make(ValidationErrors)

	if it.ProductID <= 0 {
		errs[prefix+".productId"] = "must be greater than 0"
	}
	if it.Quantity <= 0 {
		errs[prefix+".quantity"] = "must be greater than 0"
	}
	if it.PriceCents <= 0 {
		errs[prefix+".priceCents"] = "must be greater than 0"
	}
	if it.ProductTitle == "" {
		errs[prefix+".productTitle"] = "required"
	}
	if it.PriceCurr == "" {
		errs[prefix+".priceCurrency"] = "required"
	}

	return errs
}
