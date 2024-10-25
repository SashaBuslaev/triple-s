package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"triple-s/internal/config"
	"triple-s/internal/handlers"
)

func StartServer() {
	config.Parse()
	config.CreateDirAndCSV()
	http.HandleFunc("/", handlers.ListBuckets)
	http.HandleFunc("/{BucketName}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			handlers.CreateBucket(w, r)
		case "DELETE":
			handlers.DeleteBucket(w, r)
		}
	})

	fmt.Println("Server listening on port", *config.PortNum)
	PortNumStr := ":" + strconv.Itoa(*config.PortNum)
	log.Fatal(http.ListenAndServe(PortNumStr, nil))
}
