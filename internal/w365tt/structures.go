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

func (db *DbTopLevel) checkDb() {
	if db.Info.MiddayBreak == nil {
		db.Info.MiddayBreak = []int{}
	}
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
	if db.RoomGroups == nil {
		db.RoomGroups = []RoomGroup{}
	}
	if db.RoomChoiceGroups == nil {
		db.RoomChoiceGroups = []RoomChoiceGroup{}
	}
	if len(db.Classes) == 0 {
		log.Fatalln("*ERROR* No Classes")
	}
	if db.Groups == nil {
		db.Groups = []Group{}
	}
	if db.Courses == nil {
		db.Courses = []Course{}
	}
	if db.SuperCourses == nil {
		db.SuperCourses = []SuperCourse{}
	}
	if db.SubCourses == nil {
		db.SubCourses = []SubCourse{}
	}
	if db.Lessons == nil {
		db.Lessons = []Lesson{}
	}
	if db.Constraints == nil {
		db.Constraints = make(map[string]interface{})
	}
}
