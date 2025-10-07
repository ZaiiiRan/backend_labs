package dto

type V1QueryOrdersRequest struct {
	IDs               []int64 `json:"ids,omitempty"`
	CustomerIDs       []int64 `json:"customerIds,omitempty"`
	Page              int     `json:"page"`
	PageSize          int     `json:"pageSize"`
	IncludeOrderItems bool    `json:"includeOrderItems"`
}
