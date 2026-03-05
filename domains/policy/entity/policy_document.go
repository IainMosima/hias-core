package entity

import (
	"time"

	"github.com/google/uuid"
)

type PolicyDocument struct {
	ID           uuid.UUID `json:"id"`
	PolicyID     uuid.UUID `json:"policy_id"`
	MemberID     uuid.UUID `json:"member_id"`
	DocumentType string    `json:"document_type"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	S3Key        string    `json:"s3_key"`
	GeneratedBy  uuid.UUID `json:"generated_by"`
	CreatedAt    time.Time `json:"created_at"`
}
