package tool

import (
	"os"
)

func CreateFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
