package utils

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"triple-s/internal/config"
)

func ChangeBucketCSVData(bucketName string) {
	path := filepath.Join(*config.UserDir, "buckets.csv")
	records := ReadFile(path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal("Error")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	csvWriter.Write(records[0])
	records = records[1:]
	for _, record := range records {
		if record[0] == bucketName {
			record[2] = GetTime()
			record[3] = "active"
			csvWriter.Write(record)
		} else {
			csvWriter.Write(record)
		}
	}
}

func UpdateCSVObject(bucketName string, objectKey string, size int, contType string, addOrDel string) error {
	path := filepath.Join(*config.UserDir, bucketName, "objects.csv")
	records := ReadFile(path)
	records = records[1:]
	changed := false
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal("Error reading object CSV")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	csvWriter.Write([]string{"ObjectKey", "Size", "ContentType", "LastModified"})

	for _, record := range records {
		if record[0] == objectKey {
			if addOrDel == "add" {
				csvWriter.Write([]string{objectKey, strconv.Itoa(size), contType, GetTime()})
			}
			changed = true
		} else {
			csvWriter.Write(record)
		}
	}
	if !changed && addOrDel == "add" {
		csvWriter.Write([]string{objectKey, strconv.Itoa(size), contType, GetTime()})
		return nil
	} else if !changed {
		return errors.New("object not found")
	}

	return nil
}

func IsObjectPres(bucketName string, objectName string) (config.Object, bool) {
	path := filepath.Join(*config.UserDir, bucketName, "objects.csv")

	records := ReadFile(path)
	records = records[1:]
	object := config.Object{}
	for _, record := range records {
		if record[0] == objectName {
			object.Key = record[0]
			object.Size, _ = strconv.Atoi(record[1])
			object.ContentType = record[2]
			object.LastModified = record[3]
			return object, true
		}
	}
	return object, false
}
