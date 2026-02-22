package aws

import "io"

type UploadRequest struct {
	Key         string
	Body        io.Reader
	ContentType string
	Bucket      string
}

type UploadResponse struct {
	Key      string `json:"key"`
	URL      string `json:"url"`
	Bucket   string `json:"bucket"`
	Size     int64  `json:"size"`
}

type DownloadRequest struct {
	Key    string
	Bucket string
}

type DownloadResponse struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64
}
