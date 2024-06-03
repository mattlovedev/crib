package util

import (
	"encoding"
	"encoding/json"
	"log"
	"os"
)

// i is ptr to struct to be unmarshalled
func ReadFile(name string, i interface{}) error {
	/*file, err := os.Open(name)
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

	return json.Unmarshal(buf.Bytes(), i)*/
	file, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, i)
}

func BinaryReadFileNoErr(name string, i encoding.BinaryUnmarshaler) {
	file, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	if err = i.UnmarshalBinary(file); err != nil {
		log.Fatal(err)
	}
}
