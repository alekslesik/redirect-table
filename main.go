package main

import (
	"log"
	"sort"
)

type URLSet struct {
	URLs []URL `xml:"url"`
}

type URL struct {
	Loc        string `xml:"loc"`
}

func main() {
	const op = "main()"

	// Create an instance of the sample struct to pass to the function
	sampleURLSet := URLSet{}

	xmlFiles, err := getXMLFiles("./")
	if err != nil {
		log.Fatalf("%s: getXMLFiles() > %v", op, err)
	}

	xmlBytes, err := getXMLByteFiles(xmlFiles)
	if err != nil {
		log.Fatalf("%s: getXMLByteFiles() > %v", op, err)
	}

	xmlRawStructs, err := getXMLStructs(sampleURLSet, xmlBytes)
	if err != nil {
		log.Fatalf("%s: getXMLStructs() > %v", op, err)
	}

	var xmlStructs []*URLSet
	for _, xmlRawStruct := range xmlRawStructs {
		st := xmlRawStruct.(*URLSet)
		xmlStructs = append(xmlStructs, st)
	}

	// sort
	for _, s := range xmlStructs {
		sort.Slice(s.URLs, func(i, j int) bool {
			return s.URLs[i].Loc < s.URLs[j].Loc
		})
	}

	writeToCSV(xmlStructs)



}

