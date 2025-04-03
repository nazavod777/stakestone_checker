package util

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func ReadJsonFile(
	filePath string,
	fileStruct interface{},
) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error when opening file: %s", err)
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {
			log.Panicf("error when closing file: %s", err)
		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("error when reading file: %s", err)
	}

	err = json.Unmarshal(byteValue, &fileStruct)
	if err != nil {
		return fmt.Errorf("error when decoding JSON file: %s", err)
	}

	return nil
}
