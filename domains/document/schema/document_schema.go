package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/document/entity"
	"github.com/google/uuid"
)

type StandaloneDocumentResponse struct {
	ID           uuid.UUID `json:"id"`
	SourceType   string    `json:"source_type"`
	SourceID     uuid.UUID `json:"source_id"`
	DocumentType string    `json:"document_type"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	S3Key        string    `json:"s3_key"`
	CreatedBy    uuid.UUID `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type StandaloneDocumentListResponse struct {
	Documents []StandaloneDocumentResponse `json:"documents"`
	Total     int64                        `json:"total"`
	Limit     int                          `json:"limit"`
	Offset    int                          `json:"offset"`
}

type DownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
}

type UploadURLResponse struct {
	DocumentID uuid.UUID `json:"document_id"`
	UploadURL  string    `json:"upload_url"`
	S3Key      string    `json:"s3_key"`
	ExpiresIn  int64     `json:"expires_in"`
}

type DocumentResponse struct {
	ID           uuid.UUID  `json:"id"`
	EntityType   string     `json:"entity_type"`
	EntityID     uuid.UUID  `json:"entity_id"`
	DocumentType string     `json:"document_type"`
	Status       string     `json:"status"`
	FileName     string     `json:"file_name"`
	FileSize     int64      `json:"file_size"`
	MimeType     string     `json:"mime_type"`
	S3Key        string     `json:"s3_key"`
	UploadedBy   uuid.UUID  `json:"uploaded_by"`
	ConfirmedAt  *time.Time `json:"confirmed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func ToDocumentResponse(d *entity.Document) DocumentResponse {
	return DocumentResponse{
		ID:           d.ID,
		EntityType:   d.EntityType,
		EntityID:     d.EntityID,
		DocumentType: d.DocumentType,
		Status:       d.Status,
		FileName:     d.FileName,
		FileSize:     d.FileSize,
		MimeType:     d.MimeType,
		S3Key:        d.S3Key,
		UploadedBy:   d.UploadedBy,
		ConfirmedAt:  d.ConfirmedAt,
		CreatedAt:    d.CreatedAt,
	}
}
