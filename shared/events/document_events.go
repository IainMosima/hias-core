package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventDocumentUploaded       = "document.uploaded"
	EventExtractionCompleted    = "extraction.completed"
)

type DocumentUploadedEvent struct {
	DocumentID  uuid.UUID `json:"document_id"`
	UploadedBy  uuid.UUID `json:"uploaded_by"`
	FileName    string    `json:"file_name"`
	FileType    string    `json:"file_type"`
	S3Key       string    `json:"s3_key"`
	Timestamp   time.Time `json:"timestamp"`
}

type ExtractionCompletedEvent struct {
	DocumentID     uuid.UUID              `json:"document_id"`
	ExtractionData map[string]interface{} `json:"extraction_data"`
	Confidence     float64                `json:"confidence"`
	NeedsReview    bool                   `json:"needs_review"`
	Timestamp      time.Time              `json:"timestamp"`
}
