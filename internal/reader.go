package internal

import (
	"bufio"
	"fmt"
	"os"
)

func ReadData(path string) ([]string, error) {
	text := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		return text, fmt.Errorf(fmt.Sprintf("Could not open file: %s", err))
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	//read the data
	for fileScanner.Scan() {
		text = append(text, fileScanner.Text())
	}

	if err := fileScanner.Err(); err != nil {
		return text, fmt.Errorf(fmt.Sprintf("Could not scan the file: %s", err))
	}

	return text, nil
}
