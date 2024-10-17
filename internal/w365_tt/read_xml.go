package w365_tt

import (
	"encoding/xml"
	"io"
	"log"
	"os"
)

func ReadXML(filepath string) W365TT {
	// Open the  XML file
	xmlFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	// Remember to close the file at the end of the function
	defer xmlFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(xmlFile)
	log.Printf("*+ Reading: %s\n", filepath)
	v := W365TT{}
	err = xml.Unmarshal(byteValue, &v)
	if err != nil {
		log.Fatalf("XML error in %s:\n %v\n", filepath, err)
	}
	/*
	   daymap := map[string]Day{}

	   	for i, d := range v.Days {
	   		d.X = i
	   		daymap[d.Name] = d
	   	}

	   fmt.Printf("*+ Days: %+v\n", daymap)
	*/
	return v
}
