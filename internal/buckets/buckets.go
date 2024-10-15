package buckets

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

var bucketFilePath = "./buckets.csv"

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

// Validate bucket name based on S3 rules
func isValidBucketName(bucketName string) bool {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return false
	}
	// Regular expression to validate bucket name
	validNamePattern := `^[a-z0-9.-]+$`
	isValid, _ := regexp.MatchString(validNamePattern, bucketName)
	return isValid
}

// Check if the bucket name is unique
func isUniqueBucket(bucketName string) bool {
	file, err := os.Open(bucketFilePath)
	if err != nil {
		return true // Assume unique if file can't be opened
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for _, record := range records {
		if record[0] == bucketName {
			return false
		}
	}
	return true
}

// CreateBucket Handler to create a bucket
func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[len("/"):]
	if !isValidBucketName(bucketName) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}
	if !isUniqueBucket(bucketName) {
		http.Error(w, "Bucket already exists", http.StatusConflict)
		return
	}

	file, err := os.OpenFile(bucketFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	creationTime := time.Now().Format(time.RFC3339)
	writer.Write([]string{bucketName, creationTime, creationTime, "active"})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bucket %s created successfully\n", bucketName)
}

//---------------------------------------------------------------------------------------------------

// ListBuckets Handler to list all buckets
func ListBuckets(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(bucketFilePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var buckets []Bucket
	for _, record := range records {
		buckets = append(buckets, Bucket{
			Name:         record[0],
			CreationTime: record[1],
			LastModified: record[2],
			Status:       record[3],
		})
	}

	bucketList := BucketList{Buckets: buckets}
	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(bucketList)
}

//---------------------------------------------------------------------------------------

// Check if the bucket is empty
func isBucketEmpty(bucketName string) bool {
	bucketPath := fmt.Sprintf("./data/%s", bucketName)
	files, _ := ioutil.ReadDir(bucketPath)
	return len(files) == 0
}

// DeleteBucket Handler to delete a bucket
func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[len("/"):]
	if !isUniqueBucket(bucketName) { // Bucket exists
		if isBucketEmpty(bucketName) {
			deleteBucketFromCSV(bucketName)
			os.RemoveAll(fmt.Sprintf("./data/%s", bucketName))
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, "Bucket is not empty", http.StatusConflict)
		}
	} else {
		http.Error(w, "Bucket not found", http.StatusNotFound)
	}
}

// Helper function to delete bucket from CSV
func deleteBucketFromCSV(bucketName string) {
	file, err := os.Open(bucketFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	updatedRecords := [][]string{}
	for _, record := range records {
		if record[0] != bucketName {
			updatedRecords = append(updatedRecords, record)
		}
	}

	file, err = os.Create(bucketFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.WriteAll(updatedRecords)
	writer.Flush()
}
