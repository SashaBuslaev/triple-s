package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"triple-s/internal/config"
)

func ReadObjectCsv(csvName string) []config.Object {
	records := ReadFile(*config.UserDir + "/objects.csv")
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
	file, err := os.OpenFile(*config.UserDir+"/buckets.csv", os.O_RDONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	file, err = os.OpenFile(*config.UserDir+"/buckets.csv", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
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

func UpdateCSVObject(objectKey string) {
}
