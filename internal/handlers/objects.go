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
	if len(path) != 2 {
		u.CallErr(w, errors.New("invalid request"), 400)
	} else if path[1] == "" {
		u.CallErr(w, errors.New("invalid request"), 400)
	}
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
	defer file
	u.CallErr(w, err, 409)
	_, err = io.Copy(file, objectBody)
	u.CallErr(w, err, 500)
	r.Body = http.MaxBytesReader(w, r.Body, MaxSize)
	err = r.ParseMultipartForm(MaxSize)
	object := config.Object{
		Key:          objectKey,
		Size:         int(r.ContentLength),
		ContentType:  r.Header.Get("Content-Type"),
		LastModified: u.GetTime(),
	}
	u.CallErr(w, err, 400)
	u.ChangeBucketCSVData(bucketName)
	u.UpdateCSVObject(bucketName, object.Key, object.Size, object.ContentType)
	objectXML := u.GetXML(object)
	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write(objectXML)
	u.CallErr(w, err, 500)
	w.WriteHeader(http.StatusOK)
}
