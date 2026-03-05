package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClaimDocument struct {
	ID         uuid.UUID `json:"id"`
	ClaimID    uuid.UUID `json:"claim_id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	FileSize   int64     `json:"file_size"`
	S3Key      string    `json:"s3_key"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	IsDeleted  bool      `json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
