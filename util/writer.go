package util

import (
	"encoding/json"
	"os"
)

func WriteFile(name string, i interface{}) error {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	/*zipWriter := gzip.NewWriter(file)
	defer func() {
		err = zipWriter.Close()
	}()*/

	//_, err = zipWriter.Write(bytes)
	_, err = file.Write(bytes)
	return err
}
