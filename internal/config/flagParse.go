package config

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

var (
	PortNum = flag.Int("port", 8080, "Port to listen on")
	UserDir = flag.String("dir", "./data", "Directory to store user files")
)

func Parse() {
	var help = flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		PrintHelp()
		os.Exit(0)
	}

}

func PrintHelp() {
	fmt.Println("Usage:Simple Storage Service.\n\n**Usage:**\n    " +
		"triple-s [-port <N>] [-dir <S>]  \n    " +
		"triple-s --help\n\n**Options:**\n- --help     " +
		"Show this screen.\n- --port N   Port number\n- --dir S    " +
		"Path to the directory")
}

func CreateDirAndCSV() {
	if _, err := os.Stat(*UserDir); os.IsNotExist(err) {
		err := os.Mkdir(*UserDir, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	file, err := os.Create(*UserDir + "/buckets.csv")
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
