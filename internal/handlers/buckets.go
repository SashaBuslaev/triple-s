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

	err := os.Mkdir(*config.UserDir+"/"+bucketName, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.UpdateCsvBucket(bucketName)
	fmt.Println("Bucket created successfully:", bucketName)
	w.WriteHeader(http.StatusOK)

	bucket := u.GetXMLBucket(bucketName, u.GetTime(), u.GetTime(), "active")
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))
	_, err = w.Write(bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	records := u.ReadCsv(*config.UserDir + "/buckets.csv")
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return
	}
	buckets := config.BucketList{
		Buckets: records,
	}
	xmlData, err := xml.MarshalIndent(buckets, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(xmlData)
	if err != nil {
		return
	}
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {

}
