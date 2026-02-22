package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3ServiceImpl struct {
	client *s3.Client
	presignClient *s3.PresignClient
	bucket string
	cdnDomain string
}

func NewS3Service(client *s3.Client, bucket, cdnDomain string) S3Service {
	return &s3ServiceImpl{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        bucket,
		cdnDomain:     cdnDomain,
	}
}

func (s *s3ServiceImpl) Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error) {
	bucket := req.Bucket
	if bucket == "" {
		bucket = s.bucket
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(req.Key),
		Body:        req.Body,
		ContentType: aws.String(req.ContentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s/%s", s.cdnDomain, req.Key)

	return &UploadResponse{
		Key:    req.Key,
		URL:    url,
		Bucket: bucket,
	}, nil
}

func (s *s3ServiceImpl) Download(ctx context.Context, req DownloadRequest) (*DownloadResponse, error) {
	bucket := req.Bucket
	if bucket == "" {
		bucket = s.bucket
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(req.Key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	contentType := ""
	if result.ContentType != nil {
		contentType = *result.ContentType
	}

	var size int64
	if result.ContentLength != nil {
		size = *result.ContentLength
	}

	return &DownloadResponse{
		Body:        result.Body,
		ContentType: contentType,
		Size:        size,
	}, nil
}

func (s *s3ServiceImpl) GetPresignedURL(ctx context.Context, key string, expiresIn int64) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	presignResult, err := s.presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiresIn) * time.Second
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignResult.URL, nil
}

func (s *s3ServiceImpl) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}
