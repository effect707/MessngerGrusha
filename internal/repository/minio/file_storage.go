package minio

import (
	"context"
	"fmt"
	"io"

	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	client *miniogo.Client
	bucket string
}

func NewFileStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*FileStorage, error) {
	client, err := miniogo.New(endpoint, &miniogo.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &FileStorage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *FileStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, miniogo.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
	}
	return nil
}

func (s *FileStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, miniogo.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("upload file: %w", err)
	}
	return nil
}

func (s *FileStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, miniogo.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download file: %w", err)
	}
	return obj, nil
}

func (s *FileStorage) Delete(ctx context.Context, key string) error {
	if err := s.client.RemoveObject(ctx, s.bucket, key, miniogo.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
