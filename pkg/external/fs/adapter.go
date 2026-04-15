// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package fs

import (
	"context"
	"errors"
	"io"
	"strings"
)

var (
	ErrFileNotFound = errors.New("ErrFileNotFound")
)

// FileInfo holds file metadata.
type FileInfo struct {
	Name        string
	ContentType string
	Size        int64
}

func (f *FileInfo) Extension() string {
	parts := strings.Split(f.Name, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// FileData bundles a file's content reader with its metadata.
type FileData struct {
	io.Reader
	FileInfo
}

// WriteResult holds the outcome of a file write operation.
type WriteResult struct {
	// URL is a public URL to the file. File can't be accessed publicly if no URL returned
	URL string

	// Path is a path to the file that can be used for reading (works for private files that can be accessed only via Path using ReadFile method)
	Path string
}

// FileReader is an interface for reading files from a storage
type FileReader interface {
	ReadFile(ctx context.Context, path string) (*FileData, error)
}

type FileWriter interface {
	// WriteFile writes the file and returns the result with the URL or a path that can be used for reading
	// Should override the file if it already exists
	WriteFile(ctx context.Context, path string, data *FileData) (*WriteResult, error)
}

// Storage combines the FileReader and FileWriter interfaces.
type Storage interface {
	FileReader
	FileWriter
}
