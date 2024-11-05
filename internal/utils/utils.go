package utils

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"triple-s/internal/config"
)

func GetTime() string {
	now := time.Now()
	format := now.Format(time.RFC3339)
	return format
}

func CallErr(w http.ResponseWriter, err error, code int) {
	if err != nil {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(code)
		errXML := config.ErrorResponse{
			Code:    code,
			Message: err.Error(),
		}
		xmlData, _ := xml.MarshalIndent(errXML, "", "  ")
		w.Write(xmlData)
	}
}

func CreateCSVHead(header []string, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	csvFile := csv.NewWriter(file)
	defer csvFile.Flush()
	err = csvFile.Write(header)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFile(path string) [][]string {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	for _, record := range records {
		if len(record) != 4 {
			log.Fatal("CSV file parsing error")
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func GetXML(thing interface{}) []byte {
	xmlData, err := xml.MarshalIndent(thing, "", "	") // to make the text version prettier
	if err != nil {
		log.Fatal(err)
	}
	return xmlData
}

func IsValidBucket(w http.ResponseWriter, bucketName string) bool {
	if !IsValidBucketName(bucketName) {
		CallErr(w, errors.New("invalid bucket name"), 400)
		return false
	}
	if IsUniqueBucketName(bucketName) {
		CallErr(w, errors.New("bucket does not exist"), 404)
		return false
	}
	return true
}
