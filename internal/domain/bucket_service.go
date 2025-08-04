package domain

// BucketService interface for decoupling code
type BucketService interface {
	New(name string) (string, error)
	FindAllFiles(bucketName string) (*[]string, error)
	Remove(bucketName string) error
}
