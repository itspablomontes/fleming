package storage

import (
	"context"
	"io"
)

// Storage defines the interface for blob storage
type Storage interface {
	// Put uploads a blob to the storage
	Put(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error)
	// Get retrieves a blob from the storage
	Get(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	// Delete removes a blob from the storage
	Delete(ctx context.Context, bucketName, objectName string) error
	// GetURL returns a temporary URL for the object (if applicable)
	GetURL(ctx context.Context, bucketName, objectName string) (string, error)
	// CreateMultipartUpload initializes a multipart upload and returns an upload ID.
	CreateMultipartUpload(ctx context.Context, bucketName, objectName, contentType string) (string, error)
	// UploadPart uploads a single part and returns its ETag.
	UploadPart(ctx context.Context, bucketName, objectName, uploadID string, partNumber int, reader io.Reader, objectSize int64) (string, error)
	// CompleteMultipartUpload finalizes the multipart upload.
	CompleteMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string, parts []Part) (string, error)
	// AbortMultipartUpload aborts an in-progress multipart upload.
	AbortMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string) error
}

// Part represents a multipart upload part.
type Part struct {
	Number int
	ETag   string
}
