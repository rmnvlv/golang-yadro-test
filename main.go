package main

import (
	"fmt"
	"golang-yadro-test/internal"
	"log"
)

func main() {
	fmt.Println("Processing(+): Read data")
	// if len(os.Args) < 2 {
	// 	log.Fatalf("Missing parameter, provide file name!")
	// }

	data, err := internal.ReadData("./file.txt")
	if err != nil {
		log.Fatalf("Error with read data: %s", err)
	}

	fmt.Println("Processing(+): Chek data for errors and fill in normal format ")
	dataIn, err := FormatDataIn(data)
	if err != nil {
		log.Fatalf("Error with formatting data input: %s", err)
	}

	fmt.Println("Processing(+): Formatting data")
	dataOut, err := FormatDataOut(dataIn)
	if err != nil {
		log.Fatalf("Error with formatting data out: %s", err)
	}

	OutputData(dataOut)
}
