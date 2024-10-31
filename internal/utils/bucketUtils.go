package utils

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"os"
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
	r, _ := regexp.Compile("^[a-z0-9A-Z-.]+$")
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
	records := ReadCsvBucket(*config.UserDir + "/buckets.csv")
	for _, record := range records {
		fmt.Println(record)
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
		file, err := os.OpenFile(*config.UserDir+"/buckets.csv", os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		time := GetTime()
		newRecord := []string{bucketName, time, time, "Active"}
		err = writer.Write(newRecord)
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	} else if addOrDel == "del" {
		file, err := os.OpenFile(*config.UserDir+"/buckets.csv", os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		csvReader := csv.NewReader(file)
		csvWriter := csv.NewWriter(file)
		records, err := csvReader.ReadAll()
		csvWriter.Write(records[0])
		for _, record := range records {
			if record[0] != delBucket {
				csvWriter.Write(record)
			}
		}
		csvWriter.Flush()
		defer file.Close()
	}
}

func GetXMLBucket(bucketName string, creationTime string, modTime string, status string) []byte {
	bucket := config.Bucket{
		Name:         bucketName,
		CreationTime: creationTime,
		LastModified: modTime,
		Status:       status,
	}
	xmlData, err := xml.MarshalIndent(bucket, "", "	") // to make the text version prettier
	if err != nil {
		log.Fatal(err)
	}
	return xmlData
}
