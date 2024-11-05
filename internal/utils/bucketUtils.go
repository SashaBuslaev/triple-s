package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"triple-s/internal/config"
)

func IsValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}
	if name == "buckets.csv" {
		return false
	}
	if name == "internal" {
		return false
	}
	r, _ := regexp.Compile("^[a-z0-9-.]+$")
	if r.MatchString(name) {
		ipCheck, _ := regexp.Compile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if ipCheck.MatchString(name) {
			return false
		}

		if strings.Contains(name, "..") || strings.Contains(name, "--") {
			return false
		}
		if name[0] == '.' || name[0] == '-' || name[len(name)-1] == '.' || name[len(name)-1] == '-' {
			return false
		}
	}
	return true
}

func IsUniqueBucketName(name string) bool {
	path := filepath.Join(*config.UserDir, "buckets.csv")
	records := ReadCsvBucket(path)
	for _, record := range records {
		if record.Name == name {
			return false
		}
	}
	return true
}

func ReadCsvBucket(bucketName string) []config.Bucket {
	file, err := os.Open(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var buckets []config.Bucket

	csvReader := csv.NewReader(file)
	_, err = csvReader.Read() // skip first line
	if err != nil {
		return nil
	}
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		if len(record) != 4 {
			log.Fatal("ReadCSVBucket")
		}
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

func UpdateCsvBucket(bucketName string, addOrDel string, delBucket string) {
	if addOrDel == "add" {
		path := filepath.Join(*config.UserDir, "buckets.csv")
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		time := GetTime()
		newRecord := []string{bucketName, time, time, "Active"}
		err = writer.Write(newRecord)
		if err != nil {
			log.Fatal(err)
		}
	} else if addOrDel == "del" {
		path := filepath.Join(*config.UserDir, "buckets.csv")
		records := ReadFile(path)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		csvWriter := csv.NewWriter(file)
		defer csvWriter.Flush()
		csvWriter.Write(records[0])
		records = records[1:]
		for _, record := range records {
			if record[0] != delBucket {
				csvWriter.Write(record)
			}
		}
	}
}
