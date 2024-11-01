package utils

import (
	"encoding/csv"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"time"
	"triple-s/internal/config"
)

func GetTime() string {
	now := time.Now()
	format := now.Format(time.RFC3339)
	return format
}

func CallErr(w http.ResponseWriter, err error, code int) {
	if err != nil {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(code)
		errXML := config.ErrorResponse{
			Code:    code,
			Message: err.Error(),
		}
		xmlData, _ := xml.MarshalIndent(errXML, "", "  ")
		w.Write(xmlData)

	}
}

func CreateCSVHead(header []string, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	csvFile := csv.NewWriter(file)
	defer csvFile.Flush()
	err = csvFile.Write(header)
	if err != nil {
		log.Fatal(err)
	}
}
