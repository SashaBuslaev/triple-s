package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"triple-s/internal/config"
)

func CreateDirAndCSV() {
	if _, err := os.Stat(*config.UserDir); os.IsNotExist(err) {
		err := os.Mkdir(*config.UserDir, 0o777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(*config.UserDir + "/buckets.csv"); os.IsNotExist(err) {
		file, err := os.Create(*config.UserDir + "/buckets.csv")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Name", "CreationTime", "LastModified", "Status"}
		err = writer.Write(header)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Objects.csv already exists")
	}
}
