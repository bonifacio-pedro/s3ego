# S3EGO: S3 emulator built in Go
A lightweight S3 bucket emulator built in Go that allows you to create buckets, upload files, and retrieve them via simple REST endpoints or by using it as a Go module (library). Designed for local development and testing, with data stored in an in-memory SQLite database.

Remember, the project is still under development, and some tasks are yet to be added, such as:

- Unit Tests
- Swagger Docs
- In-Code Docs
- Store Metadata

## Features
- Create buckets dynamically via REST API
- Can be used as a standalone server or imported as a Go library
- Upload files to specific buckets using multipart form data
- Retrieve files by bucket and key (supports nested paths)
- Simple architecture with Go, Gin, and SQLite (in-memory)

## API Routes

| Method | Endpoint                                  | Description                      |
|--------|-------------------------------------------|---------------------------------|
| POST   | `/bucket-emulator/new-bucket/:name`       | Create a new bucket by name      |
| POST   | `/bucket-emulator/upload-file/:bucket`    | Upload a file to a bucket        |
| GET    | `/bucket-emulator/get-file/:bucket/*key`  | Download a file by key from a bucket |

## Getting Started
### Prerequisites:
- Docker installed on your machine ([Get Docker](https://docs.docker.com/get-docker/)) 
- If you use the Go Library, you need to have Go in your machine

### Run the Emulator

Simply pull and run the Docker image:

```bash
docker pull pedrobonifacio17/s3ego:latest
docker run -p 7777:7777 pedrobonifacio17/s3ego:latest
```

## CURL Examples

- Create a Bucket
```sh
curl -X POST http://localhost:7777/bucket-emulator/new-bucket/mybucket
```
- Upload a File
```sh
curl -X POST http://localhost:7777/bucket-emulator/upload-file/mybucket \
  -F "file=@/path/to/your/file.txt"
```
- Download a File
```sh
curl http://localhost:7777/bucket-emulator/get-file/mybucket/mybucket/file.txt --output downloaded_file.txt
```

## Using as a Go Library
You can also import the emulator directly in your Go projects and run it programmatically.

### Importing the package
```shell
go get "github.com/bonifacio-pedro/s3ego"
```
```go
import "github.com/bonifacio-pedro/s3ego"
```

### Starting the emulator
```go
s3 := s3ego.Start()
// s3.App gives you access to the app instance
// s3.DB is the underlying *sql.DB connection
```

### Example: Create a bucket programmatically
```go
bucket, err := s3.App.BucketService.New("mybucket")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Bucket created:", bucket.Name, bucket.Url)
```

### Functions you can use
```go
s3 := s3emulator.Start()

// Returns bucketUrl;
bucketUrl, err := s3.App.BucketService.New("mybucket")

// buckerName: string
// fileData: []byte
// fileName: string
// Returns string, error
fileKey, _ := s3.App.FileService.Upload(bucketName, fileData, fileName)

// bucketName: string
// fileKey: string
// Returns []byte, error
fileData, _ := s3.App.FileService.Get(bucketName, fileKey)

// List bucket files
// buckerName: string
// Returns []string, error
fileKeys, _ := s3.App.BucketService.FindAllFiles(bucketName)
```

## Contributing
Feel free to fork the repository, create feature branches, and submit pull requests.
