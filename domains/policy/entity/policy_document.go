package entity

import (
	"time"

	"github.com/google/uuid"
)

type PolicyDocument struct {
	ID              uuid.UUID  `json:"id"`
	PolicyID        uuid.UUID  `json:"policy_id"`
	MemberID        uuid.UUID  `json:"member_id"`
	DocumentType    string     `json:"document_type"`
	FileName        string     `json:"file_name"`
	FileSize        int64      `json:"file_size"`
	MimeType        string     `json:"mime_type"`
	S3Key           string     `json:"s3_key"`
	GeneratedBy     uuid.UUID  `json:"generated_by"`
	Version         int        `json:"version"`
	Status          string     `json:"status"`
	GenerationMode  string     `json:"generation_mode"`
	EntityType      string     `json:"entity_type"`
	EntityID        uuid.UUID  `json:"entity_id"`
	SupersededBy    *uuid.UUID `json:"superseded_by,omitempty"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	GeneratedByName string     `json:"generated_by_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
