package dto

type V1CreateAuditLogOrderRequest struct {
	Orders []V1LogOrder `json:"orders"`
}
