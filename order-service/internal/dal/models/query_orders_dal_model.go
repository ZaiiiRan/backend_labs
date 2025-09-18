package models

type QueryOrdersDalModel struct {
	Ids         []int64
	CustomerIds []int64
	Limit       int
	Offset      int
}
