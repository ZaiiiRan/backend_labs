package models

type QueryOrdersDalModel struct {
	IDs         []int64
	CustomerIDs []int64
	Limit       int
	Offset      int
}
