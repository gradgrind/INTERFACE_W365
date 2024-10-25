package readxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/db"
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func ReadXML(xmlpath string) W365TTXML {
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
	v := W365TTXML{}
	err = xml.Unmarshal(byteValue, &v)
	if err != nil {
		log.Fatalf("XML error in %s:\n %v\n", xmlpath, err)
	}
	v.Path = xmlpath
	return v
}

func ConvertToJSON(f365xml string) string {
	indata := ReadXML(f365xml)
	outdata := w365tt.W365TopLevel{}
	id2node := map[w365tt.W365Ref]interface{}{}

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
	f := strings.TrimSuffix(indata.Path, filepath.Ext(indata.Path)) + ".json"
	j, err := json.MarshalIndent(outdata, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(f, j, 0666); err != nil {
		log.Fatal(err)
	}
	return f
}

func addId(id2node map[w365tt.W365Ref]interface{}, node TTNode) w365tt.W365Ref {
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
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []Day,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		outdata.Days = append(outdata.Days, w365tt.Day{
			Id:       nid,
			Name:     n.Name,
			Shortcut: n.Shortcut,
		})
	}
}

func readHours(
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []Hour,
) {
	for i, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		r := w365tt.Hour{
			Id:       nid,
			Name:     n.Name,
			Shortcut: n.Shortcut,
		}
		t0 := get_time(n.Start)
		t1 := get_time(n.End)
		if len(t0) != 0 {
			r.Start = t0
			r.End = t1
		}
		outdata.Hours = append(outdata.Hours, r)

		if n.FirstAfternoonHour {
			outdata.W365TT.FirstAfternoonHour = i
		}
		if n.MiddayBreak {
			outdata.W365TT.MiddayBreak = append(
				outdata.W365TT.MiddayBreak, i)
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
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []Subject,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		outdata.Subjects = append(outdata.Subjects, w365tt.Subject{
			Id:       nid,
			Name:     n.Name,
			Shortcut: n.Shortcut,
		})
	}
}

func readRooms(
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []Room,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		msg := fmt.Sprintf("Room %s in RoomGroups", nid)
		rg := w365tt.GetRefList(id2node, n.RoomGroups, msg)
		if len(rg) == 0 {
			r := w365tt.Room{
				Id:       nid,
				Name:     n.Name,
				Shortcut: n.Shortcut,
			}
			msg = fmt.Sprintf("Room %s in Absences", nid)
			for _, ai := range w365tt.GetRefList(id2node, n.Absences, msg) {
				an := id2node[ai]
				r.Absences = append(r.Absences, db.TimeSlot{
					Day:  an.(Absence).Day,
					Hour: an.(Absence).Hour,
				})
			}
			sortAbsences(r.Absences)
			outdata.Rooms = append(outdata.Rooms, r)
		} else {
			r := w365tt.RoomGroup{
				Id:   nid,
				Name: n.Shortcut, // !
				//Shortcut: n.Shortcut,
				Rooms: rg,
			}
			outdata.RoomGroups = append(outdata.RoomGroups, r)
		}
	}
}

func readTeachers(
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
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
			Shortcut:         n.Shortcut,
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
		for _, ai := range w365tt.GetRefList(id2node, n.Absences, msg) {
			an := id2node[ai]
			r.Absences = append(r.Absences, db.TimeSlot{
				Day:  an.(Absence).Day,
				Hour: an.(Absence).Hour,
			})
		}
		sortAbsences(r.Absences)
		outdata.Teachers = append(outdata.Teachers, r)
	}
}

func sortAbsences(alist []db.TimeSlot) {
	slices.SortFunc(alist, func(a, b db.TimeSlot) int {
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
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []Group,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		outdata.Groups = append(outdata.Groups, w365tt.Group{
			Id:       nid,
			Shortcut: n.Shortcut,
		})
	}
}

func readClasses(
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
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
			Level:            n.Level,
			Letter:           n.Letter,
			Shortcut:         fmt.Sprintf("%d%s", n.Level, n.Letter),
			MinLessonsPerDay: n.MinLessonsPerDay,
			MaxLessonsPerDay: n.MaxLessonsPerDay,
			//MaxGapsPerDay:    -1,
			MaxAfternoons: n.MaxAfternoons,
			//MaxGapsPerWeek:   -1,
			LunchBreak:     true,
			ForceFirstHour: n.ForceFirstHour,
		}
		msg := fmt.Sprintf("Class %s in Absences", nid)
		for _, ai := range w365tt.GetRefList(id2node, n.Absences, msg) {
			an := id2node[ai]
			r.Absences = append(r.Absences, db.TimeSlot{
				Day:  an.(Absence).Day,
				Hour: an.(Absence).Hour,
			})
		}
		sortAbsences(r.Absences)
		// Initialize Divisions to get [] instead of null, when empty
		r.Divisions = []w365tt.Division{}
		msg = fmt.Sprintf("Class %s in Divisions", nid)
		for i, d := range w365tt.GetRefList(id2node, n.Divisions, msg) {
			dn := id2node[d].(Division)
			msg = fmt.Sprintf("Division %s in Groups", d)
			glist := w365tt.GetRefList(id2node, dn.Groups, msg)
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
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []Course,
) {
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		msg := fmt.Sprintf("Course %s in Subjects", nid)
		sbjs := w365tt.GetRefList(id2node, n.Subjects, msg)
		msg = fmt.Sprintf("Course %s in Groups", nid)
		grps := w365tt.GetRefList(id2node, n.Groups, msg)
		msg = fmt.Sprintf("Course %s in Teachers", nid)
		tchs := w365tt.GetRefList(id2node, n.Teachers, msg)
		msg = fmt.Sprintf("Course %s in PreferredRooms", nid)
		rms := w365tt.GetRefList(id2node, n.PreferredRooms, msg)
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
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
	items []EpochPlanCourse,
) {
	// These are currently handled as normal Courses.
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
		msg := fmt.Sprintf("EpochPlanCourse %s in Subjects", nid)
		sbjs := w365tt.GetRefList(id2node, n.Subjects, msg)
		msg = fmt.Sprintf("EpochPlanCourse %s in Groups", nid)
		grps := w365tt.GetRefList(id2node, n.Groups, msg)
		msg = fmt.Sprintf("EpochPlanCourse %s in Teachers", nid)
		tchs := w365tt.GetRefList(id2node, n.Teachers, msg)
		msg = fmt.Sprintf("EpochPlanCourse %s in PreferredRooms", nid)
		rms := w365tt.GetRefList(id2node, n.PreferredRooms, msg)
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
	outdata *w365tt.W365TopLevel,
	id2node map[w365tt.W365Ref]interface{},
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
			Id:         nid,
			Course:     n.Course,
			Duration:   dur,
			Day:        n.Day,
			Hour:       n.Hour,
			Fixed:      n.Fixed,
			LocalRooms: w365tt.GetRefList(id2node, n.LocalRooms, msg),
		})
	}
}
