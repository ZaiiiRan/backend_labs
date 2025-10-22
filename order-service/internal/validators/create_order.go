package validators

import (
	"fmt"

	pb "github.com/ZaiiiRan/backend_labs/order-service/gen/go/order-service/v1"
)

func ValidateBatchCreateRequest(req *pb.BatchCreateRequest) ValidationErrors {
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

func validateOrder(o *pb.Order, prefix string) ValidationErrors {
	errs := make(ValidationErrors)

	if o.CustomerId <= 0 {
		errs[prefix+".customer_id"] = "must be greater than 0"
	}
	if o.DeliveryAddress == "" {
		errs[prefix+".delivery_address"] = "required"
	}
	if o.TotalPriceCents <= 0 {
		errs[prefix+".total_price_cents"] = "must be greater than 0"
	}
	if o.TotalPriceCurrency == "" {
		errs[prefix+".total_price_currency"] = "required"
	}
	if len(o.OrderItems) == 0 {
		errs[prefix+".order_items"] = "at least one item is required"
		return errs
	}

	var sum int64
	currencies := map[string]struct{}{}
	for j, it := range o.OrderItems {
		iprefix := fmt.Sprintf("%s.order_items[%d]", prefix, j)
		itemErrs := validateOrderItem(it, iprefix)
		errs.Merge(itemErrs)

		sum += it.PriceCents * int64(it.Quantity)
		currencies[it.PriceCurrency] = struct{}{}
	}

	if sum != o.TotalPriceCents {
		errs[prefix+".total_price_cents"] = "must equal sum of items (priceCents * quantity)"
	}
	if len(currencies) > 1 {
		errs[prefix+".order_items.price_currency"] = "all items must have the same currency"
	}
	if len(o.OrderItems) > 0 && o.OrderItems[0].PriceCurrency != o.TotalPriceCurrency {
		errs[prefix+".total_price_currency"] = "must equal items currency"
	}

	return errs
}

func validateOrderItem(it *pb.OrderItem, prefix string) ValidationErrors {
	errs := make(ValidationErrors)

	if it.ProductId <= 0 {
		errs[prefix+".product_id"] = "must be greater than 0"
	}
	if it.Quantity <= 0 {
		errs[prefix+".quantity"] = "must be greater than 0"
	}
	if it.PriceCents <= 0 {
		errs[prefix+".price_cents"] = "must be greater than 0"
	}
	if it.ProductTitle == "" {
		errs[prefix+".product_title"] = "required"
	}
	if it.PriceCurrency == "" {
		errs[prefix+".price_currency"] = "required"
	}

	return errs
}
