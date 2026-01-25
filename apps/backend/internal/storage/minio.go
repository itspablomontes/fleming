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
	core   *minio.Core
}

func NewMinIOStorage(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	coreClient, err := minio.NewCore(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio core client: %w", err)
	}

	return &MinIOStorage{
		client: minioClient,
		core:   coreClient,
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

func (s *MinIOStorage) CreateMultipartUpload(ctx context.Context, bucketName, objectName, contentType string) (string, error) {
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	uploadID, err := s.core.NewMultipartUpload(ctx, bucketName, objectName, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to start multipart upload: %w", err)
	}

	return uploadID, nil
}

func (s *MinIOStorage) UploadPart(ctx context.Context, bucketName, objectName, uploadID string, partNumber int, reader io.Reader, objectSize int64) (string, error) {
	info, err := s.core.PutObjectPart(ctx, bucketName, objectName, uploadID, partNumber, reader, objectSize, minio.PutObjectPartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload part %d: %w", partNumber, err)
	}
	return info.ETag, nil
}

func (s *MinIOStorage) CompleteMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string, parts []Part) (string, error) {
	minioParts := make([]minio.CompletePart, 0, len(parts))
	for _, part := range parts {
		minioParts = append(minioParts, minio.CompletePart{
			ETag:       part.ETag,
			PartNumber: part.Number,
		})
	}

	_, err := s.core.CompleteMultipartUpload(ctx, bucketName, objectName, uploadID, minioParts, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to complete multipart upload: %w", err)
	}
	return objectName, nil
}

func (s *MinIOStorage) AbortMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string) error {
	if err := s.core.AbortMultipartUpload(ctx, bucketName, objectName, uploadID); err != nil {
		return fmt.Errorf("failed to abort multipart upload: %w", err)
	}
	return nil
}
