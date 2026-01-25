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
}
