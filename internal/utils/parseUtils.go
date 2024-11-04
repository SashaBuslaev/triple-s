package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"triple-s/internal/config"
)

func CreateDirAndCSV() {
	pathToCsv := filepath.Join(*config.UserDir, "buckets.csv")

	if IsExist(*config.UserDir) {
		if !IsDirEmpty(*config.UserDir) {
			if !IsExist(pathToCsv) {
				fmt.Fprintln(os.Stderr, "Choose appropriate directory with files registered in metadata")
				os.Exit(1)
			}
		}
	}
	if !IsExist(*config.UserDir) {
		fmt.Println(123)
		err := os.Mkdir(*config.UserDir, 0o777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if !IsExist(pathToCsv) {
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
	}
}

func IsDirEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	return err == io.EOF
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
