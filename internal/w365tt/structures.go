package w365tt

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
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxDays          int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
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
	Id        Ref
	Name      string
	Tag       string `json:"Shortcut"`
	Rooms     []Ref
	Reference interface{}
}

type Class struct {
	Id               Ref
	Name             string
	Tag              string `json:"Shortcut"`
	Year             int    `json:"Level"`
	Letter           string
	Absences         []TimeSlot
	Divisions        []Division
	MinLessonsPerDay interface{} `json:",omitempty"`
	MaxLessonsPerDay interface{} `json:",omitempty"`
	MaxGapsPerDay    interface{} `json:",omitempty"`
	MaxGapsPerWeek   interface{} `json:",omitempty"`
	MaxAfternoons    interface{} `json:",omitempty"`
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
	Constraints      map[string]interface{}
}
