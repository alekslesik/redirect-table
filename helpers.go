package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/fs"
	"log"
	"os"
	"reflect"
)

// getXMLFiles return slice of all xml file names in dir
func getXMLFiles(dir string) ([]string, error) {
	const op = "getXMLFiles()"

	FS := os.DirFS(dir)

	xs, err := fs.Glob(FS, "*xml")
	if err != nil {
		log.Printf("%s: glob error > %v", op, err)
		return nil, fmt.Errorf("%s: glob error > %v", op, err)
	}

	return xs, nil
}

// getXMLByteFiles return slice of converted to byte xml files
func getXMLByteFiles(files []string) ([][]byte, error) {
	const op = "getXMLByteFiles()"

	var xmlBytes [][]byte

	for _, f := range files {
		fb, err := os.ReadFile(f)
		if err != nil {
			log.Printf("%s: open file error > %v", op, err)
			return nil, fmt.Errorf("%s: open file error > %v", op, err)
		}

		xmlBytes = append(xmlBytes, fb)
	}

	return xmlBytes, nil
}

// getXMLStructs unmarshal files and returns a slice of structs
func getXMLStructs(sampleStruct interface{}, files [][]byte) ([]interface{}, error) {
	const op = "getXMLStructs()"

	var rawStructs []interface{}

	for _, xmlBytes := range files {
		// Create a new instance of the sampleStruct
		newInstance := reflect.New(reflect.TypeOf(sampleStruct)).Interface()

		// Unmarshal XML into the new instance
		err := xml.Unmarshal(xmlBytes, &newInstance)
		if err != nil {
			log.Printf("%s: unmarshal error > %v", op, err)
			return nil, fmt.Errorf("%s: unmarshal error > %v", op, err)
		}

		// Append the new instance to the result slice
		rawStructs = append(rawStructs, newInstance)
	}

	return rawStructs, nil
}

func writeToCSV(xmlStructs []*URLSet) error {
	const op = "writeToCSV()"

	l := len(xmlStructs[0].URLs)

	csvFile, err := os.OpenFile("result.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("%s: open csv file error > %v", op, err)
		return fmt.Errorf("%s: open csv file error > %v", op, err)
	}

	defer csvFile.Close()

	var rs [][]string

	for i := 0; i < l; i++ {
		a := (*xmlStructs[0]).URLs[i].Loc
		b := (*xmlStructs[1]).URLs[i].Loc

		rs = append(rs, []string{a, b})
	}


	w := csv.NewWriter(csvFile)

	for _, record := range rs {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	return nil
}
