package schema

type SendNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required,uuid"`
	Channel string `json:"channel" binding:"required,oneof=SMS EMAIL IN_APP PUSH"`
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type SendBulkNotificationRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,dive,uuid"`
	Channel string   `json:"channel" binding:"required,oneof=SMS EMAIL IN_APP PUSH"`
	Type    string   `json:"type" binding:"required"`
	Title   string   `json:"title" binding:"required"`
	Message string   `json:"message" binding:"required"`
}

type ListNotificationsRequest struct {
	UserID   string `json:"user_id" binding:"required,uuid"`
	Page     int    `json:"page" binding:"required,min=1"`
	PageSize int    `json:"page_size" binding:"required,min=1,max=100"`
}
