package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"triple-s/internal/config"
)

func ReadObjectCsv(csvName string) []config.Object {
	path := filepath.Join(*config.UserDir, "buckets.csv")
	records := ReadFile(path)
	records = records[1:]
	var objectList []config.Object
	for _, record := range records {
		if len(record) != 4 {
			log.Fatal("Wrong csv data")
		}
		object := config.Object{
			Key:          record[0],
			ContentType:  record[2],
			LastModified: record[3],
		}
		object.Size, _ = strconv.Atoi(record[1])
		objectList = append(objectList, object)
	}
	return objectList
}

func ChangeBucketCSVData(bucketName string) {
	path := filepath.Join(*config.UserDir, "buckets.csv")
	records := ReadFile(path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal("ERror")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	csvWriter.Write(records[0])
	records = records[1:]
	for _, record := range records {
		if len(record) != 4 {
			log.Fatal("Wrong csv data")
		}
		if record[0] == bucketName {
			record[2] = GetTime()
			csvWriter.Write(record)
		} else {
			csvWriter.Write(record)
		}
	}
}

func UpdateCSVObject(bucketName string, objectKey string, size int, contType string) {
	path := filepath.Join(*config.UserDir, bucketName, "objects.csv")
	records := ReadFile(path)
	records = records[1:]
	changed := false

	file, err := os.OpenFile(path, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Error reading object CSV")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	csvWriter.Write([]string{"ObjectKey", "Size", "ContentType", "LastModified"})

	for _, record := range records {
		if len(record) != 4 {
			log.Fatal("Error")
		}
		if record[0] == objectKey {
			csvWriter.Write([]string{objectKey, strconv.Itoa(size), contType, GetTime()})
			changed = true
		} else {
			csvWriter.Write(record)
		}
	}
	if !changed {
		csvWriter.Write([]string{objectKey, strconv.Itoa(size), contType, GetTime()})
	}
}
