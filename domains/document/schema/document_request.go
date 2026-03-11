package schema

type UploadURLRequest struct {
	EntityType   string `json:"entity_type" binding:"required,oneof=policy member claim quotation"`
	EntityID     string `json:"entity_id" binding:"required,uuid"`
	DocumentType string `json:"document_type" binding:"required"`
	FileName     string `json:"file_name" binding:"required"`
	FileSize     int64  `json:"file_size" binding:"required,gt=0"`
	MimeType     string `json:"mime_type" binding:"required"`
}

type BulkUploadURLRequest struct {
	Files []UploadURLRequest `json:"files" binding:"required,min=1,max=10,dive"`
}
