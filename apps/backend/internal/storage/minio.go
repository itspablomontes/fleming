package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
}

func NewMinIOStorage(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	return &MinIOStorage{
		client: minioClient,
	}, nil
}

func (s *MinIOStorage) Put(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	info, err := s.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return info.Key, nil
}

func (s *MinIOStorage) Get(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return object, nil
}

func (s *MinIOStorage) Delete(ctx context.Context, bucketName, objectName string) error {
	err := s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove object: %w", err)
	}
	return nil
}

func (s *MinIOStorage) GetURL(ctx context.Context, bucketName, objectName string) (string, error) {
	// Generate a presigned URL valid for 1 hour
	reqParams := make(map[string][]string)
	presignedURL, err := s.client.PresignedGetObject(ctx, bucketName, objectName, time.Hour, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}
