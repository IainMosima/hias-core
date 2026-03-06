package schema

import (
	"time"

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
