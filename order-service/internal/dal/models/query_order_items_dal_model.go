package models

type QueryOrderItemsDalModel struct {
	Ids      []int64
	OrderIds []int64
	Limit    int
	Offset   int
}
