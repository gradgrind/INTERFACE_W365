package w365_tt

import (
	"encoding/xml"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/base"
	"io"
	"log"
	"os"
	"strconv"
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

// TODO: Move somewhere more appropriate
func get_time(t string) string {
	// Check time and return as "mm:hh"
	tn := strings.Split(t, ":")
	if len(tn) < 2 {
		return ""
	}
	h, err := strconv.Atoi(tn[0])
	if err != nil || h > 23 || h < 0 {
		return ""
	}
	m, err := strconv.Atoi(tn[1])
	if err != nil || m > 59 || m < 0 {
		return ""
	}
	return fmt.Sprintf("%02d:%02d", h, m)
}

func collectData(w365 *W365TT, idmap IdMap) base.DBData {
	dbdata := base.NewDBData()
	add_days(&dbdata, w365.Days)
	add_hours(&dbdata, w365.Hours)
	add_teachers(&dbdata, idmap, w365.Teachers)
	return dbdata
}

func add_days(dbdata *base.DBData, items []Day) {
	for i, d := range items {
		dbdata.AddRecord(base.Record{
			"Type": "DAY", "Tag": d.Shortcut, "Name": d.Name, "X": i},
		)
	}
}

func add_hours(dbdata *base.DBData, items []Hour) {
	for i, d := range items {
		r := base.Record{
			"Type": "HOUR", "Tag": d.Shortcut, "Name": d.Name, "X": i}
		if d.FirstAfternoonHour {
			dbdata.SetInfo("AfternoonStartLesson", i)
		}
		if d.MiddayBreak {
			dbdata.AddInfo("LunchBreak", i)
		}
		t0 := get_time(d.Start)
		t1 := get_time(d.End)
		if len(t0) != 0 {
			r["StartTime"] = t0
			r["EndTime"] = t1
		}
		dbdata.AddRecord(r)
	}
}

func refList(idmap IdMap, rstring string) []interface{} {
	var reflist []interface{}
	if len(rstring) != 0 {
		for _, id := range strings.Split(rstring, ",") {
			n, ok := idmap.Id2Node[id]
			if ok {
				reflist = append(reflist, n)
			} else {
				log.Printf(" Bad Reference: %s\n", id)
			}
		}
	}
	return reflist
}

func getAbsences(idmap IdMap, alist string) [][]int {
	var absences [][]int
	for _, n0 := range refList(idmap, alist) {
		n := n0.(Absence)
		absences = append(absences, []int{n.Day, n.Hour})
	}
	return absences
}

func add_teachers(dbdata *base.DBData, idmap IdMap, items []Teacher) {
	for i, d := range items {
		r := base.Record{
			"Type": "TEACHER", "Tag": d.Shortcut, "Name": d.Name,
			"Firstname": d.Firstname, "X": i}
		absences := getAbsences(idmap, d.Absences)
		if len(absences) != 0 {
			r["NotAvailable"] = absences
		}
		/* TODO:
		MinLessonsPerDay int `xml:",attr"`
		MaxLessonsPerDay int `xml:",attr"`
		MaxDays          int `xml:",attr"`
		MaxGapsPerDay    int `xml:"MaxWindowsPerDay,attr"`
		//TODO: I have found MaxGapsPerWeek more useful
		MaxAfternoons int `xml:"NumberOfAfterNoonDays,attr"`
		*/
		dbdata.AddRecord(r)
	}
}
