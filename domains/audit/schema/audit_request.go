package schema

type ListAuditEventsRequest struct {
	Page     int `json:"page" binding:"required,min=1"`
	PageSize int `json:"page_size" binding:"required,min=1,max=100"`
}

type ListAuditEventsByEntityRequest struct {
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   string `json:"entity_id" binding:"required,uuid"`
	Page       int    `json:"page" binding:"required,min=1"`
	PageSize   int    `json:"page_size" binding:"required,min=1,max=100"`
}

type ListAuditEventsByUserRequest struct {
	UserID   string `json:"user_id" binding:"required,uuid"`
	Page     int    `json:"page" binding:"required,min=1"`
	PageSize int    `json:"page_size" binding:"required,min=1,max=100"`
}
