package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"triple-s/internal/buckets"
)

type Bucket struct {
	Name         string
	CreationTime time.Time
	LastModified time.Time
	Status       bool
}

type Object struct {
	ObjectKey    string
	Size         int64
	ContentType  string
	LastModified time.Time
}

func LandS() {
	http.HandleFunc("/", buckets.ListBuckets)
	http.HandleFunc("/{BucketName}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			buckets.CreateBucket(w, r)
		case "DELETE":
			buckets.DeleteBucket(w, r)
		}
	})
	fmt.Println(123)
	fmt.Println("Server listening on port", *PortNum)
	PortNumStr := ":" + strconv.Itoa(*PortNum)
	log.Fatal(http.ListenAndServe(PortNumStr, nil))
}
