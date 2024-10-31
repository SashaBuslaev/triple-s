package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"triple-s/internal/config"
	u "triple-s/internal/utils"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[len("/"):]

	if !u.IsValidBucketName(bucketName) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}
	if !u.IsUniqueBucketName(bucketName) {
		http.Error(w, "Bucket exists", http.StatusConflict)
		return
	}

	err := os.Mkdir(*config.UserDir+"/"+bucketName, 0o777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.UpdateCsvBucket(bucketName, "add", "")
	fmt.Println("Bucket created successfully:", bucketName)

	bucket := u.GetXMLBucket(bucketName, u.GetTime(), u.GetTime(), "active")
	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write([]byte(xml.Header))
	u.CallErr(w, err)
	_, err = w.Write(bucket)
	u.CallErr(w, err)
	w.WriteHeader(http.StatusOK)
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	records := u.ReadCsvBucket(*config.UserDir + "/buckets.csv")
	w.Header().Set("Content-Type", "application/xml")
	_, err := w.Write([]byte(xml.Header))
	u.CallErr(w, err)
	buckets := config.BucketList{
		Buckets: records,
	}
	xmlData, err := xml.MarshalIndent(buckets, "", "\t")
	u.CallErr(w, err)
	_, err = w.Write(xmlData)
	u.CallErr(w, err)
	w.WriteHeader(http.StatusOK)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketDelete := r.URL.Path[len("/"):]
	records := u.ReadCsvBucket(*config.UserDir + "/buckets.csv")
	w.Header().Set("Content-Type", "application/xml")
	bucketIs := false
	if !u.IsValidBucketName(bucketDelete) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
	}
	for _, bucketName := range records {
		if bucketName.Name == bucketDelete {
			bucketIs = true
			break
		}
	}
	if !bucketIs {
		http.Error(w, "Bucket does not exist", http.StatusNotFound)
	}
	dir, _ := os.ReadDir(bucketDelete)
	if len(dir) != 1 {
		http.Error(w, "Bucket is not empty", http.StatusConflict)
	}
	err := os.RemoveAll(*config.UserDir + "/" + bucketDelete)
	u.CallErr(w, err)
	u.UpdateCsvBucket(bucketDelete, "del", bucketDelete)
	w.WriteHeader(http.StatusNoContent)
}
