package w365tt

import (
	"gradgrind/INTERFACE_W365/internal/db"
)

// The structures used for reading a timetable-source file exported by W365.

type W365Ref string // Element reference

type Info struct {
	SchoolName         string
	Scenario           W365Ref
	Schedule           string
	FirstAfternoonHour int
	MiddayBreak        []int
}

type Day struct {
	Id       W365Ref
	Name     string
	Shortcut string
}

type Hour struct {
	Id                 W365Ref
	Name               string
	Shortcut           string
	Start              string
	End                string
	FirstAfternoonHour bool // default = false
	MiddayBreak        bool // default = false
}

type Teacher struct {
	Id               W365Ref
	Name             string
	Shortcut         string
	Firstname        string
	Absences         []db.TimeSlot
	MinLessonsPerDay interface{} `json:",omitempty"`
	MaxLessonsPerDay interface{} `json:",omitempty"`
	MaxDays          interface{} `json:",omitempty"`
	MaxGapsPerDay    interface{} `json:",omitempty"`
	MaxGapsPerWeek   interface{} `json:",omitempty"`
	MaxAfternoons    interface{} `json:",omitempty"`
	LunchBreak       bool
}

type Subject struct {
	Id       W365Ref
	Name     string
	Shortcut string
}

type Room struct {
	Id       W365Ref
	Name     string
	Shortcut string
	Absences []db.TimeSlot
}

type RoomGroup struct {
	Id       W365Ref
	Name     string
	Shortcut string
	Rooms    []W365Ref
}

type Class struct {
	Id               W365Ref
	Name             string
	Shortcut         string
	Level            int
	Letter           string
	Absences         []db.TimeSlot
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
	Id       W365Ref
	Shortcut string
}

type Division struct {
	Id     W365Ref // even though it is not top-level
	Name   string
	Groups []W365Ref
}

type Course struct {
	Id             W365Ref
	Subjects       []W365Ref `json:",omitempty"`
	Subject        W365Ref   `json:",omitempty"`
	Groups         []W365Ref
	Teachers       []W365Ref
	PreferredRooms []W365Ref
}

type SuperCourse struct {
	Id      W365Ref
	Subject W365Ref
}

type SubCourse struct {
	Id             W365Ref
	SuperCourse    W365Ref
	Subjects       []W365Ref `json:",omitempty"`
	Subject        W365Ref   `json:",omitempty"`
	Groups         []W365Ref
	Teachers       []W365Ref
	PreferredRooms []W365Ref
}

type Lesson struct {
	Id         W365Ref
	Course     W365Ref
	Duration   int
	Day        int
	Hour       int
	Fixed      bool
	LocalRooms []W365Ref
}

type W365TopLevel struct {
	W365TT       Info
	Days         []Day
	Hours        []Hour
	Teachers     []Teacher
	Subjects     []Subject
	Rooms        []Room
	RoomGroups   []RoomGroup
	Classes      []Class
	Groups       []Group
	Courses      []Course
	SuperCourses []SuperCourse
	SubCourses   []SubCourse
	Lessons      []Lesson
	Constraints  map[string]interface{}
}
