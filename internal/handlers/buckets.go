package handlers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"triple-s/internal/config"
	u "triple-s/internal/utils"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[len("/"):]

	if !u.IsValidBucketName(bucketName) {
		u.CallErr(w, errors.New("Invalid bucket name"), 400)
		return
	}
	if !u.IsUniqueBucketName(bucketName) {
		u.CallErr(w, errors.New("Bucket already exists"), 409)
		return
	}

	err := os.Mkdir(*config.UserDir+"/"+bucketName, 0o777)
	u.CallErr(w, err, 500)
	objCSVpath := path.Join(*config.UserDir+"/"+bucketName, "object.csv")
	u.CreateCSVHead([]string{"ObjectKey", "Size", "ContentType", "LastModified"}, objCSVpath)

	u.UpdateCsvBucket(bucketName, "add", "")
	fmt.Println("Bucket created successfully:", bucketName)

	bucket := u.GetXMLBucket(bucketName, u.GetTime(), u.GetTime(), "active")
	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write([]byte(xml.Header))
	u.CallErr(w, err, 500)
	_, err = w.Write(bucket)
	u.CallErr(w, err, 500)
	w.WriteHeader(http.StatusOK)
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	records := u.ReadCsvBucket(*config.UserDir + "/buckets.csv")
	w.Header().Set("Content-Type", "application/xml")
	_, err := w.Write([]byte(xml.Header))
	u.CallErr(w, err, 500)
	buckets := config.BucketList{
		Buckets: records,
	}
	xmlData, err := xml.MarshalIndent(buckets, "", "\t")
	u.CallErr(w, err, 500)
	_, err = w.Write(xmlData)
	u.CallErr(w, err, 500)
	w.WriteHeader(http.StatusOK)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketDelete := r.URL.Path[len("/"):]
	if !u.IsValidBucketName(bucketDelete) {
		u.CallErr(w, errors.New("Invalid bucket name"), 400)
	}
	records := u.ReadCsvBucket(*config.UserDir + "/buckets.csv")
	w.Header().Set("Content-Type", "application/xml")
	bucketIs := false
	for _, bucketName := range records {
		if bucketName.Name == bucketDelete {
			bucketIs = true
			break
		}
	}
	if !bucketIs {
		u.CallErr(w, errors.New("Bucket does not exist"), 404)
	}
	dir, _ := os.ReadDir(bucketDelete)
	if len(dir) != 1 {
		u.CallErr(w, errors.New("Bucket is not empty"), 409)
	}
	err := os.RemoveAll(*config.UserDir + "/" + bucketDelete)
	u.CallErr(w, err, 500)
	u.UpdateCsvBucket(bucketDelete, "del", bucketDelete)
	w.WriteHeader(http.StatusNoContent)
}
