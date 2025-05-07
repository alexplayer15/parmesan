package flags

import (
	"fmt"
	"os"
)

func DetermineOutputFileLocation(userInputOutputFileLocation string) (string, error) {
	err := validateOutputLocation(userInputOutputFileLocation)
	if err != nil {
		return "", err
	}
	os.Chdir(userInputOutputFileLocation)

	return "", err
}

func validateOutputLocation(httpFileLocation string) error {

	info, err := os.Stat(httpFileLocation)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("output location is not a directory %s", httpFileLocation)
	}

	return err
}
