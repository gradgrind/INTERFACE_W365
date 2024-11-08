package core

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// The structures used for the "database"
//TODO: Currently dealing only with the elements needed for the timetable

//+++++ These structures may be used by other database representations.

type Ref string // Element reference

type TimeSlot struct {
	Day  int
	Hour int
}

type Division struct {
	Name   string
	Groups []Ref
}

/*
type Division struct {
	Id     Ref
	Name   string
	Groups []Ref
}
*/

//------

type Info struct {
	Institution        string
	FirstAfternoonHour int
	MiddayBreak        []int
	Reference          string
}

type Day struct {
	Id   Ref
	Name string
	Tag  string
}

type Hour struct {
	Id    Ref
	Name  string
	Tag   string
	Start string
	End   string
}

type Teacher struct {
	Id               Ref
	Name             string
	Tag              string
	Firstname        string
	NotAvailable     []TimeSlot
	MinLessonsPerDay int // default = -1
	MaxLessonsPerDay int // default = -1
	MaxDays          int // default = -1
	MaxGapsPerDay    int // default = -1
	MaxGapsPerWeek   int // default = -1
	MaxAfternoons    int // default = -1
	LunchBreak       bool
}

type Subject struct {
	Id   Ref
	Name string
	Tag  string
}

type Room struct {
	Id           Ref
	Name         string
	Tag          string
	NotAvailable []TimeSlot
}

type RoomGroup struct {
	Id    Ref
	Name  string
	Tag   string
	Rooms []Ref
}

type RoomChoiceGroup struct {
	Id    Ref
	Name  string
	Tag   string
	Rooms []Ref
}

type Class struct {
	Id               Ref
	Name             string
	Tag              string
	Year             int
	Letter           string
	NotAvailable     []TimeSlot
	Divisions        []Division
	MinLessonsPerDay int // default = -1
	MaxLessonsPerDay int // default = -1
	MaxGapsPerDay    int // default = -1
	MaxGapsPerWeek   int // default = -1
	MaxAfternoons    int // default = -1
	LunchBreak       bool
	ForceFirstHour   bool
}

type Group struct {
	Id  Ref
	Tag string
}

type Course struct {
	Id       Ref
	Subject  Ref
	Groups   []Ref
	Teachers []Ref
	Room     Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type SuperCourse struct {
	Id      Ref
	Subject Ref
}

type SubCourse struct {
	Id          Ref
	SuperCourse Ref
	Subject     Ref
	Groups      []Ref
	Teachers    []Ref
	Room        Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type Lesson struct {
	Id       Ref
	Course   Ref // Course or SuperCourse Elements
	Duration int
	Day      int
	Hour     int
	Fixed    bool
	Rooms    []Ref
}

type GroupInfo struct {
	Tag      string
	Class    Ref
	Division string
}

type DbTopLevel struct {
	Info             Info
	Days             []Day
	Hours            []Hour
	Teachers         []Teacher
	Subjects         []Subject
	Rooms            []Room
	RoomGroups       []RoomGroup
	RoomChoiceGroups []RoomChoiceGroup
	Classes          []Class
	Groups           []Group
	Courses          []Course
	SuperCourses     []SuperCourse
	SubCourses       []SubCourse
	Lessons          []Lesson
	Constraints      map[string]any

	// These fields do not belong in the JSON object.
	Elements      map[Ref]any       `json:"-"`
	MaxId         int               `json:"-"` // for "indexed" Ids
	CourseLessons map[Ref][]Ref     `json:"-"`
	GroupInfoMap  map[Ref]GroupInfo `json:"-"`

	//???
	SubjectTags     map[string]Ref    `json:"-"`
	SubjectNames    map[string]string `json:"-"`
	RoomTags        map[string]Ref    `json:"-"`
	RoomChoiceNames map[string]Ref    `json:"-"`
}

func (db *DbTopLevel) NewId() Ref {
	return Ref(fmt.Sprintf("#%d", db.MaxId+1))
}

func (db *DbTopLevel) AddElement(ref Ref, element any) {
	_, nok := db.Elements[ref]
	if nok {
		Error.Fatalf("Element Id defined more than once:\n  %s\n", ref)
	}
	db.Elements[ref] = element
	// Special handling if it is an "indexed" Id.
	if strings.HasPrefix(string(ref), "#") {
		s := strings.TrimPrefix(string(ref), "#")
		i, err := strconv.Atoi(s)
		if err == nil {
			if i > db.MaxId {
				db.MaxId = i
			}
		}
	}
}

func (db *DbTopLevel) checkDb() {
	// Initializations
	if db.Info.MiddayBreak == nil {
		db.Info.MiddayBreak = []int{}
	} else {
		// Sort and check contiguity.
		slices.Sort(db.Info.MiddayBreak)
		mb := db.Info.MiddayBreak
		if mb[len(mb)-1]-mb[0] >= len(mb) {
			Error.Fatalln("MiddayBreak hours not contiguous")
		}
	}

	//??? These could be FET-specific
	db.SubjectTags = map[string]Ref{}
	db.SubjectNames = map[string]string{}
	db.RoomTags = map[string]Ref{}
	db.RoomChoiceNames = map[string]Ref{}

	// Initialize the Ref -> Element mapping
	db.Elements = make(map[Ref]any)
	if len(db.Days) == 0 {
		Error.Fatalln("No Days")
	}
	if len(db.Hours) == 0 {
		Error.Fatalln("No Hours")
	}
	if len(db.Teachers) == 0 {
		Error.Fatalln("No Teachers")
	}
	if len(db.Subjects) == 0 {
		Error.Fatalln("No Subjects")
	}
	if len(db.Rooms) == 0 {
		Error.Fatalln("No Rooms")
	}
	if len(db.Classes) == 0 {
		Error.Fatalln("No Classes")
	}
	for i, n := range db.Days {
		db.AddElement(n.Id, &db.Days[i])
	}
	for i, n := range db.Hours {
		db.AddElement(n.Id, &db.Hours[i])
	}
	for i, n := range db.Teachers {
		db.AddElement(n.Id, &db.Teachers[i])
	}
	for i, n := range db.Subjects {
		db.AddElement(n.Id, &db.Subjects[i])
	}
	for i, n := range db.Rooms {
		db.AddElement(n.Id, &db.Rooms[i])
	}
	for i, n := range db.Classes {
		db.AddElement(n.Id, &db.Classes[i])
	}
	if db.RoomGroups == nil {
		db.RoomGroups = []RoomGroup{}
	} else {
		for i, n := range db.RoomGroups {
			db.AddElement(n.Id, &db.RoomGroups[i])
		}
	}
	if db.RoomChoiceGroups == nil {
		db.RoomChoiceGroups = []RoomChoiceGroup{}
	} else {
		for i, n := range db.RoomChoiceGroups {
			db.AddElement(n.Id, &db.RoomChoiceGroups[i])
		}
	}
	if db.Groups == nil {
		db.Groups = []Group{}
	} else {
		for i, n := range db.Groups {
			db.AddElement(n.Id, &db.Groups[i])
		}
	}
	if db.Courses == nil {
		db.Courses = []Course{}
	} else {
		for i, n := range db.Courses {
			db.AddElement(n.Id, &db.Courses[i])
		}
	}
	if db.SuperCourses == nil {
		db.SuperCourses = []SuperCourse{}
	} else {
		for i, n := range db.SuperCourses {
			db.AddElement(n.Id, &db.SuperCourses[i])
		}
	}
	if db.SubCourses == nil {
		db.SubCourses = []SubCourse{}
	} else {
		for i, n := range db.SubCourses {
			db.AddElement(n.Id, &db.SubCourses[i])
		}
	}

	db.CourseLessons = map[Ref][]Ref{}
	if db.Lessons == nil {
		db.Lessons = []Lesson{}
	} else {
		for i, n := range db.Lessons {
			db.AddElement(n.Id, &db.Lessons[i])
			cref := Ref(n.Course)
			db.CourseLessons[cref] = append(db.CourseLessons[cref], n.Id)
		}
	}

	if db.Constraints == nil {
		db.Constraints = make(map[string]any)
	}

	// Get Group information
	db.GroupInfoMap = map[Ref]GroupInfo{}
	for _, c := range db.Classes {
		cid := c.Id
		for _, d := range c.Divisions {
			for _, gref := range d.Groups {
				db.GroupInfoMap[gref] = GroupInfo{
					Tag:      db.Elements[gref].(*Group).Tag,
					Class:    cid,
					Division: d.Name,
				}
			}
		}
	}
	// Check that all groups belong to a class
	for _, g := range db.Groups {
		_, ok := db.GroupInfoMap[g.Id]
		if !ok {
			delete(db.Elements, g.Id)
			Error.Printf("Group not in Class: %s\n", g.Id)
		}
	}
}

// Interface for Course and SubCourse elements
type CourseInterface interface {
	GetId() Ref
	GetGroups() []Ref
	GetTeachers() []Ref
	GetSubject() Ref
	GetRoom() Ref
}

func (c *Course) GetId() Ref            { return c.Id }
func (c *SubCourse) GetId() Ref         { return c.Id }
func (c *Course) GetGroups() []Ref      { return c.Groups }
func (c *SubCourse) GetGroups() []Ref   { return c.Groups }
func (c *Course) GetTeachers() []Ref    { return c.Teachers }
func (c *SubCourse) GetTeachers() []Ref { return c.Teachers }
func (c *Course) GetSubject() Ref       { return c.Subject }
func (c *SubCourse) GetSubject() Ref    { return c.Subject }
func (c *Course) GetRoom() Ref          { return c.Room }
func (c *SubCourse) GetRoom() Ref       { return c.Room }
