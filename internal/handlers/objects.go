package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	u "triple-s/internal/utils"
)

const MaxSize int64 = 102400

func PutObject(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/"):], "/")
	if len(path) != 2 {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	bucketName, objectKey := path[0], path[1]
	if !u.IsValidBucketName(bucketName) {
		u.CallErr(w, errors.New("Invalid bucket name"), 400)
	}
	if u.IsUniqueBucketName(bucketName) {
		u.CallErr(w, errors.New("Bucket does not exist"), 404)
	}
	objectBody := r.Body
	file, err := os.Create(objectKey)
	u.CallErr(w, err, 409)
	_, err = io.Copy(file, objectBody)
	u.CallErr(w, err, 500)
	r.Body = http.MaxBytesReader(w, r.Body, MaxSize)
	err = r.ParseMultipartForm(MaxSize)
	u.CallErr(w, err, 400)
}
