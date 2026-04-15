// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package minioFs

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"shb/pkg/external/fs"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

type MinIOStorage struct {
	fs.Storage
	client   *minio.Client
	bucket   string
	endpoint string
	useSSL   bool
	logger   *zerolog.Logger
}

type MinIOConfig struct {
	Bucket    string
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Logger    *zerolog.Logger
}

func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Ensure the target bucket exists; create it if not.
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinIOStorage{
		client:   client,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		useSSL:   cfg.UseSSL,
		logger:   cfg.Logger,
	}, nil
}

func (s *MinIOStorage) ReadFile(ctx context.Context, filePath string) (*fs.FileData, error) {
	s.logger.Debug().Str("path", filePath).Msg("MinIOStorage.ReadFile")

	obj, err := s.client.GetObject(ctx, s.bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		s.logger.Error().Err(err).Str("path", filePath).Msg("Failed to read file from MinIO")
		return nil, fs.ErrFileNotFound
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}

	return &fs.FileData{
		Reader: obj,
		FileInfo: fs.FileInfo{
			Name:        path.Base(filePath),
			ContentType: stat.ContentType,
			Size:        stat.Size,
		},
	}, nil
}

func (s *MinIOStorage) WriteFile(ctx context.Context, filePath string, data *fs.FileData) (*fs.WriteResult, error) {
	s.logger.Debug().Str("path", filePath).Msg("MinIOStorage.WriteFile")

	info, err := s.client.PutObject(ctx, s.bucket, filePath, data.Reader, data.Size, minio.PutObjectOptions{
		ContentType: data.ContentType,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("path", filePath).Msg("Failed to write file to MinIO")
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	fileURL := s.generateFileURL(filePath)

	s.logger.Debug().Str("path", filePath).Int64("size", info.Size).Msg("File written to MinIO")

	return &fs.WriteResult{
		URL:  fileURL,
		Path: filePath,
	}, nil
}

func (s *MinIOStorage) generateFileURL(filePath string) string {
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}

	// Strip the protocol prefix if the endpoint includes one.
	endpoint := s.endpoint
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		if u, err := url.Parse(endpoint); err == nil {
			endpoint = u.Host
		}
	}

	return fmt.Sprintf("%s://%s/%s/%s", scheme, endpoint, s.bucket, filePath)
}
