package handlers

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"triple-s/internal/config"
)

var bucketFilePath = "./handlers.csv"

type Bucket struct {
	Name         string `xml:"Name"`
	CreationTime string `xml:"CreationTime"`
	LastModified string `xml:"LastModified"`
	Status       string `xml:"Status"`
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[len("/"):]

	if !isValidBucketName(bucketName) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}
	if !isUniqueBucketName(bucketName) {
		http.Error(w, "Bucket exists", http.StatusConflict)
		return
	}

	err := os.Mkdir(*config.UserDir+"/"+bucketName, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Bucket created successfully:", bucketName)
	w.WriteHeader(http.StatusOK)

}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}

func isValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
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

func isUniqueBucketName(name string) bool {
	//file := os.Open(bucketFilePath)
	return true
}

func updateCsv() {}
