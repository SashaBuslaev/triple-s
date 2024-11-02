package handlers

import (
	"errors"
	"net/http"
	"strings"

	u "triple-s/internal/utils"
)

func PutObject(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[len("/"):], "/")
	if len(path) > 2 {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	bucketName, objectName := path[0], path[1]
	if !u.IsValidBucketName {
		u.CallErr(w, errors.New("Invalid bucket name"), 400)
	}
}
