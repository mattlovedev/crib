package util

import (
	"encoding"
	"encoding/json"
	"log"
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

func WriteFileNoErr(name string, i interface{}) {
	if err := WriteFile(name, i); err != nil {
		log.Fatal(err)
	}
}

func BinaryWriteFileNoErr(name string, i encoding.BinaryMarshaler) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := i.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	if _, err = file.Write(b); err != nil {
		log.Fatal(err)
	}
}
