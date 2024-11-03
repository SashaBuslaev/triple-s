package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"triple-s/internal/config"
	"triple-s/internal/handlers"

	u "triple-s/internal/utils"
)

func StartServer() {
	config.Parse()
	u.CreateDirAndCSV()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.ListBuckets(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/{BucketName}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			handlers.CreateBucket(w, r)
		case "DELETE":
			handlers.DeleteBucket(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/{BucketName}/{ObjectKey}", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[len("/"):], "/")
		if len(path) != 2 {
			u.CallErr(w, errors.New("invalid request"), 400)
		} else if path[1] == "" {
			u.CallErr(w, errors.New("invalid request"), 400)
		}
		switch r.Method {
		case "PUT":
			handlers.PutObject(w, r)
		case "GET":
			handlers.GetObject(w, r)
		case "DELETE":
			handlers.DeleteObject(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server listening on port", *config.PortNum)
	PortNumStr := ":" + strconv.Itoa(*config.PortNum)
	log.Fatal(http.ListenAndServe(PortNumStr, nil))
}
