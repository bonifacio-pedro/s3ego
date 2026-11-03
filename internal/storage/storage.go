// Package storage handles low-level file storage operations on the local filesystem,
// employing concurrent streaming pipelines and buffer pooling for high efficiency.
package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const (
	// bufferSize defines the chunk size used for streaming IO (32 KB)
	bufferSize = 32 * 1024
	// sniffSize defines the maximum bytes captured for MIME type sniffing
	sniffSize = 512
)

// FileStorage defines operations for persisting and reading raw object payloads.
type FileStorage interface {
	Store(bucketName, key string, src io.Reader) (filePath string, size int64, etag string, initialBytes []byte, err error)
	Open(filePath string) (*os.File, error)
	Delete(filePath string) error
	DeleteBucketDir(bucketName string) error
	GetBucketDir(bucketName string) string
}

type diskStorage struct {
	baseDir    string
	bufferPool *sync.Pool
}

// NewDiskStorage creates a new FileStorage instance backed by the local filesystem.
func NewDiskStorage(baseDir string) (FileStorage, error) {
	bucketsDir := filepath.Join(baseDir, "buckets")
	if err := os.MkdirAll(bucketsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create buckets directory: %w", err)
	}

	pool := &sync.Pool{
		New: func() any {
			b := make([]byte, bufferSize)
			return &b
		},
	}

	return &diskStorage{
		baseDir:    baseDir,
		bufferPool: pool,
	}, nil
}

func (s *diskStorage) GetBucketDir(bucketName string) string {
	return filepath.Join(s.baseDir, "buckets", bucketName)
}

// Store writes the incoming reader stream directly to disk while concurrently computing the MD5 ETag.
// It uses a sync.Pool for buffer allocation and goroutines for parallel disk writing and hashing.
func (s *diskStorage) Store(bucketName, key string, src io.Reader) (string, int64, string, []byte, error) {
	cleanKey := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(key, "/")), "/")
	if cleanKey == "" || cleanKey == "." {
		return "", 0, "", nil, errors.New("invalid object key")
	}

	bucketDir := s.GetBucketDir(bucketName)
	targetPath := filepath.Join(bucketDir, filepath.FromSlash(cleanKey))

	// Prevent directory traversal attacks
	cleanBucketDir := filepath.Clean(bucketDir) + string(filepath.Separator)
	cleanTargetPath := filepath.Clean(targetPath)
	if !strings.HasPrefix(cleanTargetPath+string(filepath.Separator), cleanBucketDir) && cleanTargetPath != filepath.Clean(bucketDir) {
		return "", 0, "", nil, fmt.Errorf("invalid path traversal detected for key %q", key)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", 0, "", nil, fmt.Errorf("failed to create directories for key: %w", err)
	}

	destFile, err := os.Create(targetPath)
	if err != nil {
		return "", 0, "", nil, fmt.Errorf("failed to create destination file: %w", err)
	}

	// Capture initial bytes for MIME sniffing
	bufPtr := s.bufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer s.bufferPool.Put(bufPtr)

	n, readErr := io.ReadFull(src, buf[:sniffSize])
	var initialBytes []byte
	if n > 0 {
		initialBytes = make([]byte, n)
		copy(initialBytes, buf[:n])
	}

	// Set up concurrent pipe: write to disk and compute MD5 hash in parallel
	pr, pw := io.Pipe()
	var (
		etagChan  = make(chan string, 1)
		hashErrCh = make(chan error, 1)
	)

	// Goroutine 1: Multithreaded MD5 Hashing Worker
	go func() {
		hasher := md5.New()
		hBufPtr := s.bufferPool.Get().(*[]byte)
		hBuf := *hBufPtr
		defer s.bufferPool.Put(hBufPtr)

		_, pipeErr := io.CopyBuffer(hasher, pr, hBuf)
		if pipeErr != nil {
			hashErrCh <- pipeErr
			return
		}
		etagChan <- hex.EncodeToString(hasher.Sum(nil))
	}()

	// Goroutine 2: Concurrent Disk Writer & MultiWriter to Pipe
	multiWriter := io.MultiWriter(destFile, pw)
	var written int64

	// Write captured initial bytes first
	if n > 0 {
		nw, wErr := multiWriter.Write(initialBytes)
		written += int64(nw)
		if wErr != nil {
			destFile.Close()
			pw.CloseWithError(wErr)
			os.Remove(targetPath)
			return "", 0, "", nil, fmt.Errorf("failed writing initial bytes: %w", wErr)
		}
	}

	// Stream remainder if src is not exhausted
	if readErr == nil || errors.Is(readErr, io.ErrUnexpectedEOF) {
		copied, copyErr := io.CopyBuffer(multiWriter, src, buf)
		written += copied
		if copyErr != nil && !errors.Is(copyErr, io.EOF) {
			destFile.Close()
			pw.CloseWithError(copyErr)
			os.Remove(targetPath)
			return "", 0, "", nil, fmt.Errorf("streaming copy failed: %w", copyErr)
		}
	}

	// Close stream resources
	destFile.Close()
	pw.Close()

	// Wait for hasher worker result
	select {
	case hErr := <-hashErrCh:
		os.Remove(targetPath)
		return "", 0, "", nil, fmt.Errorf("hasher error: %w", hErr)
	case etag := <-etagChan:
		return targetPath, written, etag, initialBytes, nil
	}
}

// Open opens a stored file for read operations (supporting streaming and seek).
func (s *diskStorage) Open(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

// Delete removes a single file from the disk.
func (s *diskStorage) Delete(filePath string) error {
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// DeleteBucketDir removes the entire bucket folder and all nested files.
func (s *diskStorage) DeleteBucketDir(bucketName string) error {
	dir := s.GetBucketDir(bucketName)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
