package dto

type V1QueryOrdersRequest struct {
	IDs               []int64 `json:"ids,omitempty"`
	CustomerIDs       []int64 `json:"customer_ids,omitempty"`
	Page              int     `json:"page"`
	PageSize          int     `json:"page_size"`
	IncludeOrderItems bool    `json:"include_order_items"`
}
