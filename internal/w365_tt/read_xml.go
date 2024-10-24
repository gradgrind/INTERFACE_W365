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
	v := W365TTXML{}
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

type DBItem struct {
	Id   int
	Type string
}

type IdMap struct {
	Id2Node      map[string]interface{}
	Group2Class  map[string]*Class
	Id2DBId      map[string]DBItem
	Id2RoomList  map[string][]int // RoomGroup W365Id -> rooms, list of db-ids
	Id2GroupList map[string][]int // Division W365Id -> groups, list of db-ids
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

	return IdMap{
		Id2Node:      id_node,
		Group2Class:  gid_c,
		Id2DBId:      map[string]DBItem{}, // will be populated later
		Id2RoomList:  map[string][]int{},  // will be populated later
		Id2GroupList: map[string][]int{},  // will be populated later
	}
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
	add_days(&dbdata, idmap, w365.Days)
	add_hours(&dbdata, idmap, w365.Hours)
	add_subjects(&dbdata, idmap, w365.Subjects)
	add_teachers(&dbdata, idmap, w365.Teachers)
	add_rooms(&dbdata, idmap, w365.Rooms)
	add_groups(&dbdata, idmap, w365.Groups)
	add_divisions(&dbdata, idmap, w365.Divisions)
	add_classes(&dbdata, idmap, w365.Classes)
	return dbdata
}

func add_days(dbdata *base.DBData, idmap IdMap, items []Day) {
	for i, d := range items {
		dbdata.AddRecord(base.Record{
			"Type": base.RecordType_DAY, "Tag": d.Shortcut, "Name": d.Name, "X": i},
		)
		idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_DAY}
	}
}

func add_hours(dbdata *base.DBData, idmap IdMap, items []Hour) {
	for i, d := range items {
		r := base.Record{
			"Type": base.RecordType_HOUR, "Tag": d.Shortcut, "Name": d.Name, "X": i}
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
		idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_HOUR}
	}
}

func add_subjects(dbdata *base.DBData, idmap IdMap, items []Subject) {
	for i, d := range items {
		dbdata.AddRecord(base.Record{
			"Type": base.RecordType_SUBJECT, "Tag": d.Shortcut, "Name": d.Name, "X": i},
		)
		idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_SUBJECT}
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
		n := n0.(*Absence)
		absences = append(absences, []int{n.Day, n.Hour})
	}
	return absences
}

func add_teachers(dbdata *base.DBData, idmap IdMap, items []Teacher) {
	for i, d := range items {
		r := base.Record{
			"Type": base.RecordType_TEACHER, "Tag": d.Shortcut, "Name": d.Name,
			"Firstname": d.Firstname, "X": i}
		absences := getAbsences(idmap, d.Absences)
		if len(absences) != 0 {
			r["NotAvailable"] = absences
		}
		//TODO: Is it correct to put these here?
		// They are not very hard constraints, but it might be helpful
		// to have them closely associated with the teacher.
		r["Constraints"] = map[string]int{
			"MinHoursDaily":  d.MinLessonsPerDay,
			"MaxHoursDaily":  d.MaxLessonsPerDay,
			"MaxDaysPerWeek": d.MaxDays,
			"MaxGapsPerDay":  d.MaxGapsPerDay,
			"MaxAfternoons":  d.MaxAfternoons,
			//TODO: Convert to "IntervalMaxDaysPerWeek"?
			//TODO? "MaxGapsPerWeek": d.MaxGapsPerWeek,
			//TODO? "MaxHoursDailyInInterval" for lunch break?
		}
		dbdata.AddRecord(r)
		idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_TEACHER}
	}
}

func add_rooms(dbdata *base.DBData, idmap IdMap, items []Room) {
	for i, d := range items {
		if len(d.RoomGroups) == 0 {
			r := base.Record{
				"Type": base.RecordType_ROOM, "Tag": d.Shortcut, "Name": d.Name, "X": i}
			absences := getAbsences(idmap, d.Absences)
			if len(absences) != 0 {
				r["NotAvailable"] = absences
			}
			dbdata.AddRecord(r)
			idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_ROOM}
		} else {
			var rlist []int
			for _, r := range strings.Split(d.RoomGroups, ",") {
				ritem, ok := idmap.Id2DBId[r]
				if !ok || ritem.Type != base.RecordType_ROOM {
					log.Printf(
						" *PROBLEM* Bad Room reference in RoomGroup %s:\n  %s",
						d.Id, r)
				} else {
					rlist = append(rlist, ritem.Id)
				}
			}
			idmap.Id2RoomList[d.Id] = rlist
		}
	}
}

func add_groups(dbdata *base.DBData, idmap IdMap, items []Group) {
	for i, d := range items {
		dbdata.AddRecord(base.Record{
			"Type": base.RecordType_GROUP, "Tag": d.Shortcut, "Name": d.Name, "X": i},
		)
		idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_GROUP}
	}
}

//TODO: Can these be handled later, when actually using them?
func add_divisions(dbdata *base.DBData, idmap IdMap, items []Division) {
	for _, d := range items {
		// d.Groups
		// d.Name
		var glist []int
		for _, g := range strings.Split(d.Groups, ",") {
			gitem, ok := idmap.Id2DBId[g]
			if !ok || gitem.Type != base.RecordType_GROUP {
				log.Printf(
					" *PROBLEM* Bad Group reference in RoomGroup %s:\n  %s",
					d.Id, g)
			} else {
				glist = append(glist, gitem.Id)
			}
		}
		idmap.Id2Division[d.Id] = (d.Name, glist)
	}
}

func add_classes(dbdata *base.DBData, idmap IdMap, items []Class) {
	for i, d := range items {
		r := base.Record{
			"Type": base.RecordType_CLASS, "Tag": d.Tag(), "Name": d.Name, "X": i}
		absences := getAbsences(idmap, d.Absences)
		if len(absences) != 0 {
			r["NotAvailable"] = absences
		}
		//Divisions        string
		//TODO: Can these use the raw structures?

		var plist []map[string]interface{}
		for _, p := range strings.Split(d.Divisions, ",") {
			pitem, ok := idmap.Id2GroupList[p]
			if !ok {
				log.Printf(
					" *PROBLEM* Bad Division reference in Class %s:\n  %s",
					d.Id, p)
			} else {
				plist = append(plist, map[string]interface{}{pitem)
			}
		}

		//TODO: Set field. Also correct structure.

		/* ForceFirstHour   bool
		If this is true, the following FET constraint should probably be set with
		Max_Beginnings... = 0. If false, I guess that means without the constraint.
			<ConstraintStudentsSetEarlyMaxBeginningsAtSecondHour>
				<Weight_Percentage>100</Weight_Percentage>
				<Max_Beginnings_At_Second_Hour>0</Max_Beginnings_At_Second_Hour>
				<Students>1</Students>
				<Active>true</Active>
				<Comments></Comments>
			</ConstraintStudentsSetEarlyMaxBeginningsAtSecondHour>
		*/
		//TODO: Is it correct to put these here?
		// They are not very hard constraints, but it might be helpful
		// to have them closely associated with the teacher.
		ffh := 0 // Use an int for ForceFirstHour to match the value type
		// of the other Constraints.
		if d.ForceFirstHour {
			ffh = 1
		}
		r["Constraints"] = map[string]int{
			"MinHoursDaily": d.MinLessonsPerDay,
			"MaxHoursDaily": d.MaxLessonsPerDay,
			//(TODO? "MaxGapsPerDay":  d.MaxGapsPerDay,)
			"MaxAfternoons": d.MaxAfternoons,
			//TODO: Convert to "IntervalMaxDaysPerWeek"?
			//TODO? "MaxGapsPerWeek": d.MaxGapsPerWeek,
			//TODO? "MaxHoursDailyInInterval" for lunch break?
			"ForceFirstHour": ffh,
		}
		dbdata.AddRecord(r)
		idmap.Id2DBId[d.Id] = DBItem{i, base.RecordType_CLASS}
	}
}
