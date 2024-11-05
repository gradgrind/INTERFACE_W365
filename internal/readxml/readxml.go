package readxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func ReadXML(xmlpath string) W365XML {
	// Open the  XML file
	xmlFile, err := os.Open(xmlpath)
	if err != nil {
		log.Fatal(err)
	}
	// Remember to close the file at the end of the function
	defer xmlFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(xmlFile)
	log.Printf("*+ Reading: %s\n", xmlpath)
	v := W365XML{}
	err = xml.Unmarshal(byteValue, &v)
	if err != nil {
		log.Fatalf("XML error in %s:\n %v\n", xmlpath, err)
	}
	v.Path = xmlpath
	return v
}

func ConvertToJSON(f365xml string) string {
	root := ReadXML(f365xml)
	a := root.SchoolState.ActiveScenario
	var indata *Scenario
	for i := 0; i < len(root.Scenarios); i++ {
		sp := &root.Scenarios[i]
		if sp.Id == a {
			indata = sp
			break
		}
	}
	if indata == nil {
		log.Fatalln("*ERROR* No Active Scenario")
	}
	outdata := w365tt.DbTopLevel{}
	id2node := map[w365tt.Ref]interface{}{}

	outdata.Info.Reference = string(indata.Id)
	outdata.Info.Institution = root.SchoolState.SchoolName
	//outdata.Info.Schedule = "Vorlage"
	readDays(&outdata, id2node, indata.Days)
	readHours(&outdata, id2node, indata.Hours)
	for _, n := range indata.Absences {
		id2node[n.IdStr()] = n
	}
	readSubjects(&outdata, id2node, indata.Subjects)
	readRooms(&outdata, id2node, indata.Rooms)
	readTeachers(&outdata, id2node, indata.Teachers)
	readGroups(&outdata, id2node, indata.Groups)
	for _, n := range indata.Divisions {
		id2node[n.IdStr()] = n
	}
	readClasses(&outdata, id2node, indata.Classes)
	readCourses(&outdata, id2node, indata.Courses)
	readEpochPlanCourses(&outdata, id2node, indata.EpochPlanCourses)
	readLessons(&outdata, id2node, indata.Lessons)
	// Currently no SuperCourses, SubCourses or Constraints

	// Save as JSON
	f := strings.TrimSuffix(root.Path, filepath.Ext(root.Path)) + ".json"
	j, err := json.MarshalIndent(outdata, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(f, j, 0666); err != nil {
		log.Fatal(err)
	}
	return f
}

func addId(id2node map[w365tt.Ref]interface{}, node TTNode) w365tt.Ref {
	// Check for redeclarations
	nid := node.IdStr()
	if _, ok := id2node[nid]; ok {
		log.Printf("Redefinition of %s\n", nid)
		return ""
	}
	id2node[nid] = node
	return nid
}

func readDays(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Day,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		outdata.Days = append(outdata.Days, w365tt.Day{
			Id:   nid,
			Name: n.Name,
			Tag:  n.Shortcut,
		})
	}
}

func readHours(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Hour,
) {
	for i, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		r := w365tt.Hour{
			Id:   nid,
			Name: n.Name,
			Tag:  n.Shortcut,
		}
		t0 := get_time(n.Start)
		t1 := get_time(n.End)
		if len(t0) != 0 {
			r.Start = t0
			r.End = t1
		}
		outdata.Hours = append(outdata.Hours, r)

		if n.FirstAfternoonHour {
			outdata.Info.FirstAfternoonHour = i
		}
		if n.MiddayBreak {
			outdata.Info.MiddayBreak = append(
				outdata.Info.MiddayBreak, i)
		}
	}
}

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

func readSubjects(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Subject,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		outdata.Subjects = append(outdata.Subjects, w365tt.Subject{
			Id:   nid,
			Name: n.Name,
			Tag:  n.Shortcut,
		})
	}
}

func readRooms(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Room,
) {
	rglist := map[w365tt.Ref]Room{} // RoomGroup elements
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		// Extract RoomGroup elements
		if n.RoomGroups != "" {
			rglist[nid] = n
			continue
		}
		// Normal Room
		r := w365tt.Room{
			Id:   nid,
			Name: n.Name,
			Tag:  n.Shortcut,
		}
		msg := fmt.Sprintf("Room %s in Absences", nid)
		for _, ai := range GetRefList(id2node, n.Absences, msg) {
			an := id2node[ai]
			r.NotAvailable = append(r.NotAvailable, w365tt.TimeSlot{
				Day:  an.(Absence).Day,
				Hour: an.(Absence).Hour,
			})
		}
		sortAbsences(r.NotAvailable)
		outdata.Rooms = append(outdata.Rooms, r)
	}
	// Now handle the RoomGroups
	for nid, n := range rglist {
		msg := fmt.Sprintf("Room %s in RoomGroups", nid)
		rg := GetRefList(id2node, n.RoomGroups, msg)
		r := w365tt.RoomGroup{
			Id:   nid,
			Name: n.Shortcut, // !
			//Tag: n.Shortcut,
			Rooms: rg,
		}
		outdata.RoomGroups = append(outdata.RoomGroups, r)
	}
}

func readTeachers(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Teacher,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		r := w365tt.Teacher{
			Id:               nid,
			Name:             n.Name,
			Tag:              n.Shortcut,
			Firstname:        n.Firstname,
			MinLessonsPerDay: n.MinLessonsPerDay,
			MaxLessonsPerDay: n.MaxLessonsPerDay,
			MaxDays:          n.MaxDays,
			MaxGapsPerDay:    n.MaxGapsPerDay,
			//MaxGapsPerWeek:   -1,
			MaxAfternoons: n.MaxAfternoons,
			LunchBreak:    true,
		}
		msg := fmt.Sprintf("Teacher %s in Absences", nid)
		for _, ai := range GetRefList(id2node, n.Absences, msg) {
			an := id2node[ai]
			r.NotAvailable = append(r.NotAvailable, w365tt.TimeSlot{
				Day:  an.(Absence).Day,
				Hour: an.(Absence).Hour,
			})
		}
		sortAbsences(r.NotAvailable)
		outdata.Teachers = append(outdata.Teachers, r)
	}
}

func sortAbsences(alist []w365tt.TimeSlot) {
	slices.SortFunc(alist, func(a, b w365tt.TimeSlot) int {
		if a.Day < b.Day {
			return -1
		}
		if a.Day == b.Day {
			if a.Hour < b.Hour {
				return -1
			}
			if a.Hour == b.Hour {
				log.Fatalln("Equal Absences")
			}
			return 1
		}
		return 1
	})
}

func readGroups(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Group,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		outdata.Groups = append(outdata.Groups, w365tt.Group{
			Id:  nid,
			Tag: n.Shortcut,
		})
	}
}

func readClasses(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Class,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		r := w365tt.Class{
			Id:               nid,
			Name:             n.Name,
			Year:             n.Level,
			Letter:           n.Letter,
			Tag:              fmt.Sprintf("%d%s", n.Level, n.Letter),
			MinLessonsPerDay: n.MinLessonsPerDay,
			MaxLessonsPerDay: n.MaxLessonsPerDay,
			//MaxGapsPerDay:    -1,
			MaxAfternoons: n.MaxAfternoons,
			//MaxGapsPerWeek:   -1,
			LunchBreak:     true,
			ForceFirstHour: n.ForceFirstHour,
		}
		msg := fmt.Sprintf("Class %s in Absences", nid)
		for _, ai := range GetRefList(id2node, n.Absences, msg) {
			an := id2node[ai]
			r.NotAvailable = append(r.NotAvailable, w365tt.TimeSlot{
				Day:  an.(Absence).Day,
				Hour: an.(Absence).Hour,
			})
		}
		sortAbsences(r.NotAvailable)
		// Initialize Divisions to get [] instead of null, when empty
		r.Divisions = []w365tt.Division{}
		msg = fmt.Sprintf("Class %s in Divisions", nid)
		for i, d := range GetRefList(id2node, n.Divisions, msg) {
			dn := id2node[d].(Division)
			msg = fmt.Sprintf("Division %s in Groups", d)
			glist := GetRefList(id2node, dn.Groups, msg)
			if len(glist) != 0 {
				nm := dn.Name
				if nm == "" {
					nm = fmt.Sprintf("#div%d", i)
				}
				r.Divisions = append(r.Divisions, w365tt.Division{
					Name:   nm,
					Groups: glist,
				})
			}
		}
		outdata.Classes = append(outdata.Classes, r)
	}
}

func readCourses(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Course,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		msg := fmt.Sprintf("Course %s in Subjects", nid)
		sbjs := GetRefList(id2node, n.Subjects, msg)
		msg = fmt.Sprintf("Course %s in Groups", nid)
		grps := GetRefList(id2node, n.Groups, msg)
		msg = fmt.Sprintf("Course %s in Teachers", nid)
		tchs := GetRefList(id2node, n.Teachers, msg)
		msg = fmt.Sprintf("Course %s in PreferredRooms", nid)
		rms := GetRefList(id2node, n.PreferredRooms, msg)
		outdata.Courses = append(outdata.Courses, w365tt.Course{
			Id:             nid,
			Subjects:       sbjs,
			Groups:         grps,
			Teachers:       tchs,
			PreferredRooms: rms,
		})
	}
}

func readEpochPlanCourses(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []EpochPlanCourse,
) {
	// These are currently handled as normal Courses.
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		msg := fmt.Sprintf("EpochPlanCourse %s in Subjects", nid)
		sbjs := GetRefList(id2node, n.Subjects, msg)
		msg = fmt.Sprintf("EpochPlanCourse %s in Groups", nid)
		grps := GetRefList(id2node, n.Groups, msg)
		msg = fmt.Sprintf("EpochPlanCourse %s in Teachers", nid)
		tchs := GetRefList(id2node, n.Teachers, msg)
		msg = fmt.Sprintf("EpochPlanCourse %s in PreferredRooms", nid)
		rms := GetRefList(id2node, n.PreferredRooms, msg)
		outdata.Courses = append(outdata.Courses, w365tt.Course{
			Id:             nid,
			Subjects:       sbjs,
			Groups:         grps,
			Teachers:       tchs,
			PreferredRooms: rms,
		})
	}
}

func readLessons(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Lesson,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		if _, ok := id2node[n.Course]; !ok {
			log.Printf("Lesson with invalid Course: %s\n", nid)
			continue
		}
		dur := 1
		if n.DoubleLesson {
			// Currently none, pending changes
			dur = 2
		}
		msg := fmt.Sprintf("Course %s in LocalRooms", nid)
		outdata.Lessons = append(outdata.Lessons, w365tt.Lesson{
			Id:       nid,
			Course:   n.Course,
			Duration: dur,
			Day:      n.Day,
			Hour:     n.Hour,
			Fixed:    n.Fixed,
			Rooms:    GetRefList(id2node, n.LocalRooms, msg),
		})
	}
}

func GetRefList(
	id2node map[w365tt.Ref]interface{},
	reflist RefList,
	messages ...string,
) []w365tt.Ref {
	var rl []w365tt.Ref
	if reflist != "" {
		for _, rs := range strings.Split(string(reflist), ",") {
			rr := w365tt.Ref(rs)
			if _, ok := id2node[rr]; ok {
				rl = append(rl, rr)
			} else {
				log.Printf("Invalid Reference in RefList: %s\n", rs)
				for _, msg := range messages {
					log.Printf("  ++ %s\n", msg)
				}
			}
		}
	}
	return rl
}
