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
