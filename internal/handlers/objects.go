package handlers

import (
	"errors"
	"fmt"
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
	if !u.IsValidBucket(w, bucketName) {
		return
	}
	fmt.Println("haha")
	if objectKey == "objects.csv" {
		u.CallErr(w, errors.New("forbidden object name"), http.StatusBadRequest)
		return
	}
	objectBody := r.Body
	objectPath := filepath.Join(*config.UserDir, bucketName, objectKey)
	file, err := os.Create(objectPath)
	if u.CallErr(w, err, 409) {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, objectBody)
	if u.CallErr(w, err, 500) {
		return
	}
	if r.ContentLength >= MaxSize {
		u.CallErr(w, errors.New("max limit reached, upgrade your subscription plan to get more than 100 mb"), 400)
		return
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
	if u.CallErr(w, err, 500) {
		return
	}
	// w.WriteHeader(http.StatusOK)
}

func GetObject(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/"):], "/")
	bucketName, objectKey := path[0], path[1]
	if !u.IsValidBucket(w, bucketName) {
		return
	}

	object, isPres := u.IsObjectPres(bucketName, objectKey)

	if !isPres {
		u.CallErr(w, errors.New("object not found"), 404)
		return
	}

	w.Header().Set("Content-Type", object.ContentType)
	objectBody, err := io.ReadAll(r.Body)
	if u.CallErr(w, err, 500) {
		return
	}
	w.Write(objectBody)
	// w.WriteHeader(http.StatusOK)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/"):], "/")
	bucketName, objectKey := path[0], path[1]
	pathToObj := filepath.Join(*config.UserDir, bucketName, objectKey)
	if !u.IsValidBucket(w, bucketName) {
		return
	}

	isPres := u.UpdateCSVObject(bucketName, objectKey, 0, "", "del")
	if u.CallErr(w, isPres, 404) {
		return
	}

	err := os.Remove(pathToObj)
	if u.CallErr(w, err, 404) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
	u.ChangeBucketCSVData(bucketName)
}
