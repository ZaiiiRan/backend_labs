package models

type QueryOrderItemsDalModel struct {
	IDs      []int64
	OrderIDs []int64
	Limit    int
	Offset   int
}
