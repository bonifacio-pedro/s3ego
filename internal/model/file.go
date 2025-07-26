package model

import "fmt"

type File struct {
	ID       int    `json:"id"`
	Key      string `json:"key"`
	Data     []byte `json:"data"`
	BucketID uint   `json:"bucket_id"`
}

func NewFile(data []byte, bucket Bucket, fileName string) File {
	file := new(File)
	file.Data = data
	file.BucketID = uint(bucket.ID)
	file.Key = fmt.Sprintf("%s/%s", bucket.Name, fileName)

	return *file
}
