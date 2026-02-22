package aws

import "context"

type S3Service interface {
	Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error)
	Download(ctx context.Context, req DownloadRequest) (*DownloadResponse, error)
	GetPresignedURL(ctx context.Context, key string, expiresIn int64) (string, error)
	Delete(ctx context.Context, key string) error
}
