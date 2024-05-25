package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
)

// i is ptr to struct to be unmarshalled
func ReadFile(name string, i interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	zipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() {
		err = zipReader.Close()
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(zipReader)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf.Bytes(), i)
}
