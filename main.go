package main

import (
	"flag"
	"fmt"
	"golang-yadro-test/internal"
	"log"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "./test_file.txt", "Provide the file path")
	flag.Parse()

	fmt.Println("Processing(+): Read data")
	data, err := internal.ReadData(path)
	if err != nil {
		log.Fatalf("Error with read data: %s", err)
	}

	fmt.Println("Processing(+): Chek data for errors and fill in normal format ")
	formattedData, err := internal.FormatData(data)
	if err != nil {
		log.Fatalf("Error with formatting data input: %s", err)
	}

	fmt.Println("Processing(+): Formatting data")
	parsedData, err := internal.ParseData(formattedData)
	if err != nil {
		log.Fatalf("Error with formatting data out: %s", err)
	}

	internal.WriteData(parsedData)
}
