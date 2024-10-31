package utils

import (
	"net/http"
	"time"
)

func GetTime() string {
	now := time.Now()
	format := now.Format(time.RFC3339)
	return format
}

func CallErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
