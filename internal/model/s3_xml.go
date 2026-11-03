package model

import (
	"encoding/xml"
	"time"
)

// ListBucketResult represents the standard AWS S3 XML response for GET /{bucket}
type ListBucketResult struct {
	XMLName        xml.Name        `xml:"ListBucketResult"`
	Xmlns          string          `xml:"xmlns,attr"`
	Name           string          `xml:"Name"`
	Prefix         string          `xml:"Prefix"`
	KeyCount       int             `xml:"KeyCount"`
	MaxKeys        int             `xml:"MaxKeys"`
	IsTruncated    bool            `xml:"IsTruncated"`
	Contents       []ObjectSummary `xml:"Contents"`
	CommonPrefixes []string        `xml:"CommonPrefixes>Prefix,omitempty"`
}

// ObjectSummary represents an object item inside ListBucketResult
type ObjectSummary struct {
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int64     `xml:"Size"`
	StorageClass string    `xml:"StorageClass"`
}

// S3ErrorResponse represents the standard AWS S3 XML error payload
type S3ErrorResponse struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	Resource  string   `xml:"Resource,omitempty"`
	RequestID string   `xml:"RequestId,omitempty"`
}
