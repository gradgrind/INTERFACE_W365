package w365_tt

import (
	"encoding/xml"
	"gradgrind/INTERFACE_W365/internal/base"
	"io"
	"log"
	"os"
	"strings"
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

type IdMap struct {
	Id2Node     map[string]interface{}
	Group2Class map[string]*Class
}

func makeIdMap(w365 *W365TT) IdMap {
	id_node := map[string]interface{}{}

	for i := 0; i < len(w365.Days); i++ {
		n := &(w365.Days[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Hours); i++ {
		n := &(w365.Hours[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Absences); i++ {
		n := &(w365.Absences[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Teachers); i++ {
		n := &(w365.Teachers[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Subjects); i++ {
		n := &(w365.Subjects[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Rooms); i++ {
		n := &(w365.Rooms[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Classes); i++ {
		n := &(w365.Classes[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Groups); i++ {
		n := &(w365.Groups[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Divisions); i++ {
		n := &(w365.Divisions[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Courses); i++ {
		n := &(w365.Courses[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.EpochPlanCourses); i++ {
		n := &(w365.EpochPlanCourses[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Lessons); i++ {
		n := &(w365.Lessons[i])
		id_node[n.IdStr()] = n
	}
	for i := 0; i < len(w365.Fractions); i++ {
		n := &(w365.Fractions[i])
		id_node[n.IdStr()] = n
	}

	gid_c := map[string]*Class{}
	for i := 0; i < len(w365.Classes); i++ {
		c := &(w365.Classes[i])
		for _, gid := range strings.Split(c.Groups, ",") {
			gid_c[gid] = c
		}
	}

	return IdMap{id_node, gid_c}
}

func collectData(w365 *W365TT, idmap IdMap) base.DBData {
	dbdata := base.NewDBData()
	for i, d := range w365.Days {
		dbdata.AddRecord(base.Record{
			"Type": "DAY", "Tag": d.Shortcut, "Name": d.Name, "X": i},
		)

	}
	return dbdata
}
