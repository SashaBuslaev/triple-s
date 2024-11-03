package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"triple-s/internal/config"

	u "triple-s/internal/utils"
)

const MaxSize int64 = 102400000

func PutObject(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/"):], "/")
	bucketName, objectKey := path[0], path[1]
	if !u.IsValidBucketName(bucketName) {
		u.CallErr(w, errors.New("invalid bucket name"), 400)
	}
	if u.IsUniqueBucketName(bucketName) {
		u.CallErr(w, errors.New("bucket does not exist"), 404)
	}
	objectBody := r.Body
	objectPath := filepath.Join(*config.UserDir, bucketName, objectKey)
	file, err := os.Create(objectPath)
	u.CallErr(w, err, 409)
	defer file.Close()
	_, err = io.Copy(file, objectBody)
	u.CallErr(w, err, 500)
	if r.ContentLength >= MaxSize {
		u.CallErr(w, errors.New("max limit reached, upgrade your subscription plan to get more than 100 mb"), 400)
	}
	r.Body = http.MaxBytesReader(w, r.Body, MaxSize) // if misleading file length
	object := config.Object{
		Key:          objectKey,
		Size:         int(r.ContentLength),
		ContentType:  r.Header.Get("Content-Type"),
		LastModified: u.GetTime(),
	}

	u.ChangeBucketCSVData(bucketName)
	u.UpdateCSVObject(bucketName, object.Key, object.Size, object.ContentType, "add")
	objectXML := u.GetXML(object)
	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write(objectXML)
	u.CallErr(w, err, 500)
	w.WriteHeader(http.StatusOK)
}

func GetObject(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/"):], "/")
	bucketName, objectKey := path[0], path[1]
	if !u.IsValidBucketName(bucketName) {
		u.CallErr(w, errors.New("invalid bucket name"), 400)
	}
	if u.IsUniqueBucketName(bucketName) {
		u.CallErr(w, errors.New("bucket does not exist"), 404)
	}
	if !u.IsObjectPres(bucketName, objectKey) {
		u.CallErr(w, errors.New("object not found"), 404)
	}
	w.Header().Set("Content-Type", r.)
}
