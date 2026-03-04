package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type QuotationDocument struct {
	ID             uuid.UUID       `json:"id"`
	QuotationID    uuid.UUID       `json:"quotation_id"`
	VersionNumber  int             `json:"version_number"`
	FileName       string          `json:"file_name"`
	FileType       string          `json:"file_type"`
	FileSize       int64           `json:"file_size"`
	S3Key          string          `json:"s3_key"`
	UploadedBy     uuid.UUID       `json:"uploaded_by"`
	CanEditRoles   json.RawMessage `json:"can_edit_roles"`
	CanDeleteRoles json.RawMessage `json:"can_delete_roles"`
	IsDeleted      bool            `json:"is_deleted"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
