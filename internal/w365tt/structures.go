package w365tt

import "log"

// The structures used for the "database", adapted to read from W365
//TODO: Currently dealing only with the elements needed for the timetable

type Ref string // Element reference

type Info struct {
	Institution        string `json:"SchoolName"`
	FirstAfternoonHour int
	MiddayBreak        []int
	Reference          string `json:"Scenario"`
}

type Day struct {
	Id   Ref
	Name string
	Tag  string `json:"Shortcut"`
}

type Hour struct {
	Id    Ref
	Name  string
	Tag   string `json:"Shortcut"`
	Start string
	End   string
	// These are for W365 only, optional, with default = False:
	FirstAfternoonHour bool `json:",omitempty"`
	MiddayBreak        bool `json:",omitempty"`
}

type TimeSlot struct {
	Day  int
	Hour int
}

type Teacher struct {
	Id               Ref
	Name             string
	Tag              string `json:"Shortcut"`
	Firstname        string
	NotAvailable     []TimeSlot `json:"Absences"`
	MinLessonsPerDay interface{}
	MaxLessonsPerDay interface{}
	MaxDays          interface{}
	MaxGapsPerDay    interface{}
	MaxGapsPerWeek   interface{}
	MaxAfternoons    interface{}
	LunchBreak       bool
}

type Subject struct {
	Id   Ref
	Name string
	Tag  string `json:"Shortcut"`
}

type Room struct {
	Id           Ref
	Name         string
	Tag          string     `json:"Shortcut"`
	NotAvailable []TimeSlot `json:"Absences"`
}

type RoomGroup struct {
	Id    Ref
	Name  string
	Tag   string `json:"Shortcut"`
	Rooms []Ref
}

type RoomChoiceGroup struct {
	Id    Ref
	Name  string
	Tag   string `json:"Shortcut"`
	Rooms []Ref
}

type Class struct {
	Id               Ref
	Name             string
	Tag              string `json:"Shortcut"`
	Year             int    `json:"Level"`
	Letter           string
	NotAvailable     []TimeSlot `json:"Absences"`
	Divisions        []Division
	MinLessonsPerDay interface{}
	MaxLessonsPerDay interface{}
	MaxGapsPerDay    interface{}
	MaxGapsPerWeek   interface{}
	MaxAfternoons    interface{}
	LunchBreak       bool
	ForceFirstHour   bool
}

type Group struct {
	Id  Ref
	Tag string `json:"Shortcut"`
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

type Course struct {
	Id             Ref
	Subjects       []Ref `json:",omitempty"`
	Subject        Ref
	Groups         []Ref
	Teachers       []Ref
	PreferredRooms []Ref `json:",omitempty"`
	// Not in W365:
	Room Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type SuperCourse struct {
	Id      Ref
	Subject Ref
}

type SubCourse struct {
	Id             Ref
	SuperCourse    Ref
	Subjects       []Ref `json:",omitempty"`
	Subject        Ref
	Groups         []Ref
	Teachers       []Ref
	PreferredRooms []Ref `json:",omitempty"`
	// Not in W365:
	Room Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type Lesson struct {
	Id       Ref
	Course   Ref // Course or SuperCourse Elements
	Duration int
	Day      int
	Hour     int
	Fixed    bool
	Rooms    []Ref `json:"LocalRooms"` // only Room Elements
}

type DbTopLevel struct {
	Info             Info `json:"W365TT"`
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
	Constraints      map[string]interface{}
}

func (db *DbTopLevel) checkDb() map[Ref]interface{} {
	if db.Info.MiddayBreak == nil {
		db.Info.MiddayBreak = []int{}
	}
	// Initialize the Ref -> Element mapping
	dbrefs := make(map[Ref]interface{})
	if len(db.Days) == 0 {
		log.Fatalln("*ERROR* No Days")
	}
	if len(db.Hours) == 0 {
		log.Fatalln("*ERROR* No Hours")
	}
	if len(db.Teachers) == 0 {
		log.Fatalln("*ERROR* No Teachers")
	}
	if len(db.Subjects) == 0 {
		log.Fatalln("*ERROR* No Subjects")
	}
	if len(db.Rooms) == 0 {
		log.Fatalln("*ERROR* No Rooms")
	}
	if len(db.Classes) == 0 {
		log.Fatalln("*ERROR* No Classes")
	}
	for i, n := range db.Days {
		addId(dbrefs, n.Id, &db.Days[i])
	}
	for i, n := range db.Hours {
		addId(dbrefs, n.Id, &db.Hours[i])
	}
	for i, n := range db.Teachers {
		addId(dbrefs, n.Id, &db.Teachers[i])
	}
	for i, n := range db.Subjects {
		addId(dbrefs, n.Id, &db.Subjects[i])
	}
	for i, n := range db.Rooms {
		addId(dbrefs, n.Id, &db.Rooms[i])
	}
	for i, n := range db.Classes {
		addId(dbrefs, n.Id, &db.Classes[i])
	}
	if db.RoomGroups == nil {
		db.RoomGroups = []RoomGroup{}
	} else {
		for i, n := range db.RoomGroups {
			addId(dbrefs, n.Id, &db.RoomGroups[i])
		}
	}
	if db.RoomChoiceGroups == nil {
		db.RoomChoiceGroups = []RoomChoiceGroup{}
	} else {
		for i, n := range db.RoomChoiceGroups {
			addId(dbrefs, n.Id, &db.RoomChoiceGroups[i])
		}
	}
	if db.Groups == nil {
		db.Groups = []Group{}
	} else {
		for i, n := range db.Groups {
			addId(dbrefs, n.Id, &db.Groups[i])
		}
	}
	if db.Courses == nil {
		db.Courses = []Course{}
	} else {
		for i, n := range db.Courses {
			addId(dbrefs, n.Id, &db.Courses[i])
		}
	}
	if db.SuperCourses == nil {
		db.SuperCourses = []SuperCourse{}
	} else {
		for i, n := range db.SuperCourses {
			addId(dbrefs, n.Id, &db.SuperCourses[i])
		}
	}
	if db.SubCourses == nil {
		db.SubCourses = []SubCourse{}
	} else {
		for i, n := range db.SubCourses {
			addId(dbrefs, n.Id, &db.SubCourses[i])
		}
	}
	if db.Lessons == nil {
		db.Lessons = []Lesson{}
	} else {
		for i, n := range db.Lessons {
			addId(dbrefs, n.Id, &db.Lessons[i])
		}
	}
	if db.Constraints == nil {
		db.Constraints = make(map[string]interface{})
	}
	return dbrefs
}

func addId(refs map[Ref]interface{}, id Ref, node interface{}) {
	_, nok := refs[id]
	if nok {
		log.Fatalf("*ERROR* Element Id defined more than once:\n  %s\n", id)
	}
	refs[id] = node
}
