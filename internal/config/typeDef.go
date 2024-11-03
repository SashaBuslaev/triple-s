package config

import "encoding/xml"

type Bucket struct {
	Name         string `xml:"Name"`
	CreationTime string `xml:"CreationTime"`
	LastModified string `xml:"LastModified"`
	Status       string `xml:"Status"`
}

type BucketList struct {
	XMLName xml.Name `xml:"Buckets"`
	Buckets []Bucket `xml:"Bucket"`
}

type ErrorResponse struct {
	XMLName xml.Name `xml:"ErrorResponse"`
	Code    int      `xml:"Code"`
	Message string   `xml:"Message"`
}

type Object struct {
	Key          string `xml:"Name"`
	Size         int    `xml:"Size in bytes"`
	ContentType  string `xml:"Type"`
	LastModified string `xml:"LastModified"`
}
