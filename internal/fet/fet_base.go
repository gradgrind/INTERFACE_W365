// Package fet handles interaction with the fet timetabling program.
package fet

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"log"
	"strings"
)

type Ref = w365tt.Ref

const fet_version = "6.25.2"

// Function makeXML produces a chunk of pretty-printed XML output from
// the input data.
func makeXML(data interface{}, indent_level int) string {
	const indent = "  "
	prefix := strings.Repeat(indent, indent_level)
	xmlData, err := xml.MarshalIndent(data, prefix, indent)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	return string(xmlData)
}

type fet struct {
	Version          string `xml:"version,attr"`
	Mode             string
	Institution_Name string
	Comments         string // this can be a source reference
	Days_List        fetDaysList
	Hours_List       fetHoursList
	Teachers_List    fetTeachersList
	Subjects_List    fetSubjectsList
	Rooms_List       fetRoomsList
	Students_List    fetStudentsList
	//Buildings_List
	Activity_Tags_List     fetActivityTags
	Activities_List        fetActivitiesList
	Time_Constraints_List  timeConstraints
	Space_Constraints_List spaceConstraints
}

type virtualRoom struct {
	rooms       []Ref   // only Rooms
	roomChoices [][]Ref // list of pure Room lists
}

type courseInfo struct {
	subject  Ref
	groups   []Ref
	teachers []Ref
	room     virtualRoom
	lessons  []Ref
}

type fetInfo struct {
	db            *w365tt.DbTopLevel
	ref2fet       map[Ref]string
	ref2grouponly map[Ref]string
	days          []string
	hours         []string
	fetdata       fet
	// These cover only courses and groups with lessons:
	superSubs       map[Ref][]Ref
	courseInfo      map[Ref]courseInfo // Key can be Course or SuperCourse
	classDivisions  map[Ref][][]Ref
	atomicGroups    map[Ref][]AtomicGroup
	fetVirtualRooms map[string]string

	//TODO ... ???
	//courses          map[Ref]int
	//subcourses       map[Ref]int
	//supercourses     map[Ref]int
	//fixed_activities []bool
}

type timeConstraints struct {
	XMLName xml.Name `xml:"Time_Constraints_List"`
	//
	ConstraintBasicCompulsoryTime                basicTimeConstraint
	ConstraintStudentsSetNotAvailableTimes       []studentsNotAvailable
	ConstraintTeacherNotAvailableTimes           []teacherNotAvailable
	ConstraintActivityPreferredStartingTime      []startingTime
	ConstraintMinDaysBetweenActivities           []minDaysBetweenActivities
	ConstraintStudentsSetMaxHoursDailyInInterval []lunchBreak
	ConstraintStudentsSetMaxGapsPerWeek          []maxGapsPerWeek
	ConstraintStudentsSetMinHoursDaily           []minLessonsPerDay
}

type basicTimeConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsoryTime"`
	Weight_Percentage int
	Active            bool
}

type spaceConstraints struct {
	XMLName                          xml.Name `xml:"Space_Constraints_List"`
	ConstraintBasicCompulsorySpace   basicSpaceConstraint
	ConstraintActivityPreferredRoom  []fixedRoom
	ConstraintActivityPreferredRooms []roomChoice
}

type basicSpaceConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsorySpace"`
	Weight_Percentage int
	Active            bool
}

func make_fet_file(dbdata *w365tt.DbTopLevel,

// activities []wzbase.Activity,
// course2activities map[int][]int,
// subject_activities []wzbase.SubjectGroupActivities,
) string {
	//TODO--
	//fmt.Printf("\n????? %+v\n", dbdata.Info)

	// Build ref-index -> fet-key mapping
	ref2fet := map[Ref]string{}
	for _, r := range dbdata.Subjects {
		ref2fet[r.Id] = r.Tag
	}
	for _, r := range dbdata.Rooms {
		ref2fet[r.Id] = r.Tag
	}
	for _, r := range dbdata.Teachers {
		ref2fet[r.Id] = r.Tag
	}
	ref2grouponly := map[Ref]string{}
	for _, r := range dbdata.Groups {
		ref2grouponly[r.Id] = r.Tag
	}
	for _, r := range dbdata.Classes {
		ref2fet[r.Id] = r.Tag
		// Handle the groups
		for _, d := range r.Divisions {
			for _, g := range d.Groups {
				ref2fet[g] = r.Tag + CLASS_GROUP_SEP + ref2grouponly[g]
			}
		}
	}

	//fmt.Printf("ref2fet: %v\n", ref2fet)

	fetinfo := fetInfo{
		db:            dbdata,
		ref2fet:       ref2fet,
		ref2grouponly: ref2grouponly,
		fetdata: fet{
			Version:          fet_version,
			Mode:             "Official",
			Institution_Name: dbdata.Info.Institution,
			Comments:         getString(dbdata.Info.Reference),
			Time_Constraints_List: timeConstraints{
				ConstraintBasicCompulsoryTime: basicTimeConstraint{
					Weight_Percentage: 100, Active: true},
			},
			Space_Constraints_List: spaceConstraints{
				ConstraintBasicCompulsorySpace: basicSpaceConstraint{
					Weight_Percentage: 100, Active: true},
			},
		},
		fetVirtualRooms: map[string]string{},
	}

	getDays(&fetinfo)
	getHours(&fetinfo)
	getTeachers(&fetinfo)
	getSubjects(&fetinfo)
	getRooms(&fetinfo)
	fmt.Println("=====================================")
	gatherCourseInfo(&fetinfo)

	//readCourseIndexes(&fetinfo)
	makeAtomicGroups(&fetinfo)
	//fmt.Println("\n +++++++++++++++++++++++++++")
	//printAtomicGroups(&fetinfo)
	getClasses(&fetinfo)
	getActivities(&fetinfo)

	/*
		gap_subject_activities(&fetinfo, subject_activities)
	*/

	return xml.Header + makeXML(fetinfo.fetdata, 0)
}

func getString(val interface{}) string {
	s, ok := val.(string)
	if !ok {
		b, _ := json.Marshal(val)
		s = string(b)
	}
	return s
}
