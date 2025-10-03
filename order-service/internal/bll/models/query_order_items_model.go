package models

type QueryOrderItemsModel struct {
	IDs               []int64
	CustomerIDs       []int64
	Page              int
	PageSize          int
	IncludeOrderItems bool
}
