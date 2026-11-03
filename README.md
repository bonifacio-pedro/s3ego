# S3EGO: Lightweight S3 Wire-Compatible Emulator in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/bonifacio-pedro/s3ego#section-documentation.svg)](https://pkg.go.dev/github.com/bonifacio-pedro/s3ego#section-documentation)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bonifacio-pedro/s3ego)](https://golang.org)
[![Docker](https://img.shields.io/docker/pulls/pedrobonifacio17/s3ego)](https://hub.docker.com/r/pedrobonifacio17/s3ego)

A high-performance, lightweight **Amazon S3 Wire-Compatible** (Path-Style) bucket emulator built in Go. Designed for local development, CI/CD pipelines, and automated integration testing without requiring cloud credentials or heavy infrastructure.

Data is stored persistently on the filesystem (MinIO-style) paired with a pure-Go SQLite metadata catalog (`modernc.org/sqlite`, CGO-free), utilizing streaming pipelines and concurrent goroutines for high-throughput uploads and low memory consumption.

---

## Highlights

- **S3 Wire Protocol (Path-Style):** Direct HTTP/REST compatibility with AWS S3 XML schemas (`ListBucketResult`, `<Error>`).
- **Binary Stream Uploads:** Upload objects via raw HTTP body (`PUT /{bucket}/{key...}`) without multipart/form-data overhead.
- **Concurrent Streaming Pipeline:** Employs `sync.Pool` buffer reuse and multithreaded goroutines to write chunks to disk and calculate MD5 `ETag` simultaneously.
- **MinIO-Style Persistence:** Raw payloads reside directly on disk (`S3EGO_DATA_DIR`), while object metadata is indexed in a lightweight, pure-Go SQLite database.
- **Dual Mode:** Run as a standalone Docker container or import programmatically as a Go library with automatic isolated temporary storage for unit tests.
- **Permissive Auth (Mock Mode):** Transparently accepts any AWS SigV4 signed requests or unsigned anonymous requests. Works out-of-the-box with AWS CLI, `boto3`, and official AWS SDKs.
- **Zero-CGO Static Binary:** Generates small static binaries (`CGO_ENABLED=0`) suited for `alpine` and `scratch` container images.

---

## Benchmark: S3EGO vs MinIO

An empirical benchmark was conducted under identical conditions comparing **S3EGO** against **MinIO** (`minio/minio:RELEASE.2024-01-16T16-07-38Z`) using containerized environments and standard AWS SigV4 authenticated requests over concurrent workers:

| Metric / Workload | S3EGO | MinIO | Relative Advantage |
| :--- | :--- | :--- | :--- |
| **Docker Image Size** | **38.8 MB** (10.8 MB comp.) | 207.0 MB (52.9 MB comp.) | **5.3x smaller image** |
| **Idle Memory (RAM)** | **2.15 MiB** | 93.4 MiB | **43.4x less RAM** |
| **Peak Memory (Under Load)** | **13.21 MiB** | 245.4 MiB | **18.6x less RAM** |
| **Cold Start / Boot Time** | **< 50 ms** | ~1.8 s | **36x faster boot** |
| **Large Upload (5 MB $\times$ 20, $C=4$)** | **389.61 MB/s** | 338.42 MB/s | **+15.1% higher throughput** |
| **Large Upload Latency** | **44.67 ms** | 57.48 ms | **22.3% lower latency** |
| **Small Upload (20 KB $\times$ 300, $C=15$)** | **308.9 req/s** | 300.0 req/s | **+3.0% higher RPS** |
| **List Objects (`ListObjectsV2` XML)** | **310.0 req/s** | 131.9 req/s | **2.35x faster (3.2ms vs 7.5ms)** |
| **Large Download (5 MB)** | **1698.90 MB/s** | 1642.82 MB/s | **Comparable** (loopback saturated) |
| **Small Download (20 KB)** | 4595.1 req/s | **7590.0 req/s** | MinIO internal cache advantage |

### Key Benchmark Insights
- **18.6x to 43x Less Memory:** While MinIO initializes clustering meshes, web consoles, and Prometheus exporters taking ~100MB idle and ~250MB under load, S3EGO operates between **2 MiB and 14 MiB**, enabling dozens of test instances without host exhaustion.
- **Lock-Free Concurrent Streaming:** By decoupling disk writing and MD5 `ETag` calculation into asynchronous goroutines communicating over `io.Pipe`, S3EGO outpaced MinIO by **15.1% on large uploads** (389.6 MB/s vs 338.4 MB/s).
- **Fast B-Tree Metadata Indexing:** `ListObjectsV2` queries run in **3.23 ms** (vs 7.58 ms on MinIO) thanks to direct indexed SQLite range scans (`idx_files_bucket_key`) without filesystem tree crawling.

> 📄 **Full Technical Paper:** For comprehensive methodology, mathematical equations, and in-depth architectural comparisons, see [result.tex](result.tex).

---

## S3 Path-Style API Endpoints

| HTTP Method | Endpoint | S3 Action | Description |
| :--- | :--- | :--- | :--- |
| `PUT` | `/{bucket}` | `CreateBucket` | Creates a new S3 bucket (200 OK) |
| `HEAD` | `/{bucket}` | `HeadBucket` | Verifies bucket existence (200 OK or 404) |
| `GET` | `/{bucket}` | `ListObjectsV2` | Returns S3 standard XML `<ListBucketResult>` |
| `DELETE` | `/{bucket}` | `DeleteBucket` | Deletes a bucket and its stored files (204 No Content) |
| `PUT` | `/{bucket}/{key...}` | `PutObject` | Uploads binary stream as an object, returns `ETag` |
| `GET` | `/{bucket}/{key...}` | `GetObject` | Downloads object with streaming and range support |
| `HEAD` | `/{bucket}/{key...}` | `HeadObject` | Returns object headers (`ETag`, `Content-Length`, `Content-Type`) |
| `DELETE` | `/{bucket}/{key...}` | `DeleteObject` | Deletes an object (204 No Content) |

---

## Quickstart with Docker

### Run with Docker CLI
```bash
docker run -d \
  --name s3ego \
  -p 7777:7777 \
  -v $(pwd)/data:/data \
  pedrobonifacio17/s3ego:latest
```

### Run with Docker Compose
```bash
docker compose up -d
```

---

## Usage Examples

### 1. Using `curl`

- **Create a Bucket**
  ```bash
  curl -X PUT http://localhost:7777/my-bucket
  ```

- **Upload an Object (Raw Binary Payload)**
  ```bash
  curl -X PUT http://localhost:7777/my-bucket/nested/folder/photo.jpg \
    -H "Content-Type: image/jpeg" \
    --data-binary "@/path/to/local/photo.jpg"
  ```

- **Download an Object**
  ```bash
  curl -OJ http://localhost:7777/my-bucket/nested/folder/photo.jpg
  ```

- **List Objects in Bucket (S3 XML)**
  ```bash
  curl http://localhost:7777/my-bucket
  # Or with prefix:
  curl "http://localhost:7777/my-bucket?prefix=nested/"
  ```

- **Delete an Object**
  ```bash
  curl -X DELETE http://localhost:7777/my-bucket/nested/folder/photo.jpg
  ```

- **Delete a Bucket**
  ```bash
  curl -X DELETE http://localhost:7777/my-bucket
  ```

---

### 2. Using AWS CLI

Configure a dummy profile (or pass mock credentials) pointing to S3EGO:

```bash
# List objects in bucket
aws --endpoint-url=http://localhost:7777 s3 ls s3://my-bucket/

# Copy file to bucket
aws --endpoint-url=http://localhost:7777 s3 cp ./document.pdf s3://my-bucket/docs/document.pdf

# Download file from bucket
aws --endpoint-url=http://localhost:7777 s3 cp s3://my-bucket/docs/document.pdf ./downloaded.pdf

# Sync folder
aws --endpoint-url=http://localhost:7777 s3 sync ./my-folder s3://my-bucket/backup/
```

---

### 3. Using as an Embedded Go Library

You can embed S3EGO directly into your Go test suites:

```bash
go get github.com/bonifacio-pedro/s3ego
```

```go
package main_test

import (
    "bytes"
    "io"
    "testing"

    "github.com/bonifacio-pedro/s3ego"
)

func TestMyStorageService(t *testing.T) {
    // Starts an isolated in-memory/temp S3EGO instance
    s3 := s3ego.Start()
    defer s3.Close()

    // Create a bucket
    err := s3.Bucket.Create("test-bucket")
    if err != nil {
        t.Fatal(err)
    }

    // Upload an object
    payload := []byte("Hello from Go test!")
    fileMeta, err := s3.File.Put("test-bucket", "docs/hello.txt", bytes.NewReader(payload))
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("Uploaded ETag: %s, Size: %d", fileMeta.ETag, fileMeta.Size)

    // Retrieve the object
    meta, fileReader, err := s3.File.Get("test-bucket", "docs/hello.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer fileReader.Close()

    data, _ := io.ReadAll(fileReader)
    t.Logf("Read content: %s", string(data))
}
```

To persist data to a specific folder instead of a temp directory, pass the `WithDataDir` functional option:
```go
s3 := s3ego.Start(s3ego.WithDataDir("./test-storage"))
defer s3.Close()
```

---

## Environment Variables

| Variable | Default | Description |
| :--- | :--- | :--- |
| `S3EGO_DATA_DIR` | `./data` | Directory where raw objects (`./data/buckets/...`) and metadata DB (`./data/metadata.db`) are persisted |
| `S3EGO_PORT` | `7777` | HTTP port on which the S3 emulator listens |

---

## Architecture & Storage Engine

```
       +---------------------------------------------+
       | HTTP Request (PUT/GET /{bucket}/{key...})   |
       +---------------------------------------------+
                              |
                     [ Gin Router & Middleware ]
                      - S3 Headers & CORS
                      - Permissive Auth Bypass
                              |
       +----------------------+----------------------+
       |                                             |
[ File Storage Engine ]                     [ SQLite Pure-Go ]
 - sync.Pool Buffer Allocation               - modernc.org/sqlite (CGO-free)
 - Goroutine 1: Multithreaded MD5 ETag       - Fast B-Tree Indexes
 - Goroutine 2: Concurrent Disk Streaming    - Bucket & Object Metadatas
 - Direct Filesystem (./data/buckets/...)    - Persistent across restarts
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
