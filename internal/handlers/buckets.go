package handlers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/internal/config"

	u "triple-s/internal/utils"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[len("/"):]

	if !u.IsValidBucketName(bucketName) {
		u.CallErr(w, errors.New("invalid bucket name"), 400)
		return
	}
	if !u.IsUniqueBucketName(bucketName) {
		u.CallErr(w, errors.New("bucket already exists"), 409)
		return
	}
	path := filepath.Join(*config.UserDir, bucketName)
	err := os.Mkdir(path, 0o777)
	if u.CallErr(w, err, 500) {
		return
	}
	objCSVpath := filepath.Join(path, "objects.csv")
	u.CreateCSVHead([]string{"ObjectKey", "Size", "ContentType", "LastModified"}, objCSVpath)

	u.UpdateCsvBucket(bucketName, "add", "")
	fmt.Println("Bucket created successfully:", bucketName)

	bucket := config.Bucket{
		Name:         bucketName,
		CreationTime: u.GetTime(),
		LastModified: u.GetTime(),
		Status:       "inactive",
	}

	bucketXML := u.GetXML(bucket)
	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write([]byte(xml.Header))
	if u.CallErr(w, err, 500) {
		return
	}
	_, err = w.Write(bucketXML)
	if u.CallErr(w, err, 500) {
		return
	}
	// w.WriteHeader(http.StatusOK)
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(*config.UserDir, "buckets.csv")
	records := u.ReadCsvBucket(path)
	w.Header().Set("Content-Type", "application/xml")
	_, err := w.Write([]byte(xml.Header))
	if u.CallErr(w, err, 500) {
		return
	}
	buckets := config.BucketList{
		Buckets: records,
	}
	xmlData, err := xml.MarshalIndent(buckets, "", "\t")
	if u.CallErr(w, err, 500) {
		return
	}
	_, err = w.Write(xmlData)
	if u.CallErr(w, err, 500) {
		return
	}
	// w.WriteHeader(http.StatusOK)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketDelete := r.URL.Path[len("/"):]
	if !u.IsValidBucketName(bucketDelete) {
		u.CallErr(w, errors.New("invalid bucket name"), 400)
		return
	}
	path := filepath.Join(*config.UserDir, "buckets.csv")
	records := u.ReadCsvBucket(path)
	w.Header().Set("Content-Type", "application/xml")
	bucketIs := false
	for _, bucketName := range records {
		if bucketName.Name == bucketDelete {
			bucketIs = true
			break
		}
	}
	if !bucketIs {
		u.CallErr(w, errors.New("bucket does not exist"), 404)
		return
	}
	path = filepath.Join(*config.UserDir, bucketDelete)
	dir, _ := os.ReadDir(path)
	if len(dir) != 1 {
		fmt.Println(len(dir))
		u.CallErr(w, errors.New("bucket is not empty"), 409)
		return
	}
	path = filepath.Join(*config.UserDir, bucketDelete)
	err := os.RemoveAll(path)
	if u.CallErr(w, err, 500) {
		return
	}
	u.UpdateCsvBucket(bucketDelete, "del", bucketDelete)
	w.WriteHeader(http.StatusNoContent)
}
