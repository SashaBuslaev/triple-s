package utils

import (
	"encoding/csv"
	"log"
	"os"
	"time"
	"triple-s/internal/config"
)

func ReadCsv(bucketName string) []config.Bucket {
	file, err := os.Open(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var buckets []config.Bucket

	csvReader := csv.NewReader(file)
	_, err = csvReader.Read() //skip first line
	if err != nil {
		return nil
	}
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		bucket := config.Bucket{
			Name:         record[0],
			CreationTime: record[1],
			LastModified: record[2],
			Status:       record[3],
		}
		buckets = append(buckets, bucket)
	}
	return buckets
}

func GetTime() string {
	now := time.Now()
	format := now.Format(time.RFC3339)
	return format
}
