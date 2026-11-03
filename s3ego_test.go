package s3ego_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bonifacio-pedro/s3ego"
	"github.com/bonifacio-pedro/s3ego/internal/model"
)

func TestS3WireCompatibleProtocol(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "s3ego-e2e-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	s3 := s3ego.Start(s3ego.WithDataDir(tempDir))
	defer s3.Close()

	router := s3.App.Router
	bucketName := "e2e-bucket"

	// 1. Create Bucket: PUT /{bucket}
	t.Run("CreateBucket", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/"+bucketName, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		if loc := w.Header().Get("Location"); loc != "/"+bucketName {
			t.Errorf("expected Location header %q, got %q", "/"+bucketName, loc)
		}
		if s3Server := w.Header().Get("Server"); s3Server != "AmazonS3" {
			t.Errorf("expected Server header AmazonS3, got %q", s3Server)
		}
	})

	// 2. Head Bucket: HEAD /{bucket}
	t.Run("HeadBucket", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/"+bucketName, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		// Non-existent bucket
		req404 := httptest.NewRequest(http.MethodHead, "/non-existent-bucket", nil)
		w404 := httptest.NewRecorder()
		router.ServeHTTP(w404, req404)
		if w404.Code != http.StatusNotFound {
			t.Errorf("expected status 404 for missing bucket, got %d", w404.Code)
		}
	})

	// 3. Put Object: PUT /{bucket}/{key...} with binary body stream
	testContent := []byte("Hello, this is a test payload for S3EGO wire compatibility!")
	hasher := md5.New()
	hasher.Write(testContent)
	expectedETag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hasher.Sum(nil)))

	t.Run("PutObject", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/"+bucketName+"/nested/folder/file.txt", bytes.NewReader(testContent))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		etag := w.Header().Get("ETag")
		if etag != expectedETag {
			t.Errorf("expected ETag %s, got %s", expectedETag, etag)
		}
	})

	// 4. Head Object: HEAD /{bucket}/{key...}
	t.Run("HeadObject", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/"+bucketName+"/nested/folder/file.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if etag := w.Header().Get("ETag"); etag != expectedETag {
			t.Errorf("expected ETag %s, got %s", expectedETag, etag)
		}
		if cl := w.Header().Get("Content-Length"); cl != fmt.Sprintf("%d", len(testContent)) {
			t.Errorf("expected Content-Length %d, got %s", len(testContent), cl)
		}
	})

	// 5. Get Object: GET /{bucket}/{key...}
	t.Run("GetObject", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+bucketName+"/nested/folder/file.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if etag := w.Header().Get("ETag"); etag != expectedETag {
			t.Errorf("expected ETag %s, got %s", expectedETag, etag)
		}
		if !bytes.Equal(w.Body.Bytes(), testContent) {
			t.Errorf("expected body %q, got %q", string(testContent), w.Body.String())
		}
	})

	// 6. List Objects (XML): GET /{bucket}
	t.Run("ListObjectsXML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+bucketName, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var result model.ListBucketResult
		if err := xml.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to parse XML ListBucketResult: %s", err)
		}

		if result.Name != bucketName {
			t.Errorf("expected bucket name %q, got %q", bucketName, result.Name)
		}
		if len(result.Contents) != 1 {
			t.Fatalf("expected 1 object in contents, got %d", len(result.Contents))
		}
		if result.Contents[0].Key != "nested/folder/file.txt" {
			t.Errorf("expected key 'nested/folder/file.txt', got %q", result.Contents[0].Key)
		}
		if result.Contents[0].ETag != expectedETag {
			t.Errorf("expected ETag %s, got %s", expectedETag, result.Contents[0].ETag)
		}
	})

	// 7. Delete Object: DELETE /{bucket}/{key...}
	t.Run("DeleteObject", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/"+bucketName+"/nested/folder/file.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d", w.Code)
		}

		// Verify deletion
		getReq := httptest.NewRequest(http.MethodGet, "/"+bucketName+"/nested/folder/file.txt", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)
		if getW.Code != http.StatusNotFound {
			t.Errorf("expected status 404 after deletion, got %d", getW.Code)
		}
	})

	// 8. Delete Bucket: DELETE /{bucket}
	t.Run("DeleteBucket", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/"+bucketName, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d", w.Code)
		}

		// Verify bucket deletion
		headReq := httptest.NewRequest(http.MethodHead, "/"+bucketName, nil)
		headW := httptest.NewRecorder()
		router.ServeHTTP(headW, headReq)
		if headW.Code != http.StatusNotFound {
			t.Errorf("expected status 404 after bucket deletion, got %d", headW.Code)
		}
	})
}

func TestPersistenceAcrossRestarts(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "s3ego-persist-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	payload := []byte("Persistent data that survives restart")

	// Phase 1: Start S3EGO, write bucket and object, then shutdown
	{
		s3 := s3ego.Start(s3ego.WithDataDir(tempDir))
		err := s3.Bucket.Create("persistent-bucket")
		if err != nil {
			t.Fatalf("failed creating bucket in phase 1: %s", err)
		}

		_, err = s3.File.Put("persistent-bucket", "data/persistent.txt", bytes.NewReader(payload))
		if err != nil {
			t.Fatalf("failed storing file in phase 1: %s", err)
		}
		s3.Close()
	}

	// Phase 2: Start new S3EGO instance pointing to the same data directory
	{
		s3New := s3ego.Start(s3ego.WithDataDir(tempDir))
		defer s3New.Close()

		exists, err := s3New.Bucket.Exists("persistent-bucket")
		if err != nil || !exists {
			t.Fatalf("expected bucket to persist across restart, exists=%v err=%v", exists, err)
		}

		meta, fileHandle, err := s3New.File.Get("persistent-bucket", "data/persistent.txt")
		if err != nil {
			t.Fatalf("failed to retrieve file after restart: %s", err)
		}
		defer fileHandle.Close()

		readBytes, err := io.ReadAll(fileHandle)
		if err != nil {
			t.Fatalf("failed reading retrieved file: %s", err)
		}

		if !bytes.Equal(readBytes, payload) {
			t.Errorf("persisted payload mismatch: expected %q, got %q", string(payload), string(readBytes))
		}
		if meta.Size != int64(len(payload)) {
			t.Errorf("expected size %d, got %d", len(payload), meta.Size)
		}
	}
}

func TestConcurrentLargeFileStreaming(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "s3ego-concurrent-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	s3 := s3ego.Start(s3ego.WithDataDir(tempDir))
	defer s3.Close()

	if err := s3.Bucket.Create("large-bucket"); err != nil {
		t.Fatalf("failed creating bucket: %s", err)
	}

	// Generate 1 MB deterministic payload
	chunk := bytes.Repeat([]byte("0123456789ABCDEF"), 64) // 1 KB
	largePayload := bytes.Repeat(chunk, 1024)               // 1 MB

	hasher := md5.New()
	hasher.Write(largePayload)
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	req := httptest.NewRequest(http.MethodPut, "/large-bucket/large.bin", bytes.NewReader(largePayload))
	w := httptest.NewRecorder()
	s3.App.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	expectedETag := fmt.Sprintf(`"%s"`, expectedHash)
	if etag := w.Header().Get("ETag"); etag != expectedETag {
		t.Errorf("expected ETag %s, got %s", expectedETag, etag)
	}

	// Retrieve and verify data integrity
	meta, fileHandle, err := s3.File.Get("large-bucket", "large.bin")
	if err != nil {
		t.Fatalf("failed to retrieve stored large object: %s", err)
	}
	defer fileHandle.Close()

	if meta.Size != int64(len(largePayload)) {
		t.Errorf("expected size %d, got %d", len(largePayload), meta.Size)
	}

	retrieved, err := io.ReadAll(fileHandle)
	if err != nil {
		t.Fatalf("failed reading object body: %s", err)
	}

	if !bytes.Equal(retrieved, largePayload) {
		t.Errorf("data corruption: retrieved bytes differ from sent bytes")
	}
}

func TestS3ErrorResponses(t *testing.T) {
	s3 := s3ego.Start()
	defer s3.Close()

	router := s3.App.Router

	// 1. Get object from nonexistent bucket
	req := httptest.NewRequest(http.MethodGet, "/unknown-bucket/some-key.txt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	var errResp model.S3ErrorResponse
	if err := xml.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse S3 XML error: %s", err)
	}
	if errResp.Code != "NoSuchBucket" {
		t.Errorf("expected error code NoSuchBucket, got %s", errResp.Code)
	}

	// 2. Create bucket with invalid name
	reqBad := httptest.NewRequest(http.MethodPut, "/INVALID_UPPERCASE_NAME", nil)
	wBad := httptest.NewRecorder()
	router.ServeHTTP(wBad, reqBad)

	if wBad.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid bucket name, got %d", wBad.Code)
	}
	var badResp model.S3ErrorResponse
	if err := xml.Unmarshal(wBad.Body.Bytes(), &badResp); err != nil {
		t.Fatalf("failed to parse S3 XML error: %s", err)
	}
	if badResp.Code != "InvalidBucketName" {
		t.Errorf("expected error code InvalidBucketName, got %s", badResp.Code)
	}
}
