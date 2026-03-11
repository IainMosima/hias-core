package entity

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID
	EntityType   string
	EntityID     uuid.UUID
	DocumentType string
	Status       string
	FileName     string
	FileSize     int64
	MimeType     string
	S3Key        string
	UploadedBy   uuid.UUID
	ConfirmedAt  *time.Time
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
