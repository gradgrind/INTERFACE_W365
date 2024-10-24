package w365tt

import (
	"gradgrind/INTERFACE_W365/internal/db"
	"log"
	"strings"
)

// The structures used for reading a timetable-source file exported by W365.

type W365Ref string     // Element reference
type W365RefList string // "List" of Element references

func GetRefList(
	id2node map[W365Ref]interface{},
	reflist W365RefList,
	messages ...string,
) []W365Ref {
	var rl []W365Ref
	if reflist != "" {
		for _, rs := range strings.Split(string(reflist), ",") {
			rr := W365Ref(rs)
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

type TTNode interface {
	// An interface for Top-Level-Elements with Id field
	IdStr() W365Ref
}

type Info struct {
	FirstAfternoonHour int
	MiddayBreak        []int
}

type Day struct {
	Id       W365Ref
	Type     string
	Name     string
	Shortcut string
}

func (n *Day) IdStr() W365Ref {
	return n.Id
}

type Hour struct {
	Id                 W365Ref
	Type               string
	Name               string
	Shortcut           string
	Start              string
	End                string
	FirstAfternoonHour bool // default = false
	MiddayBreak        bool // default = false
}

func (n *Hour) IdStr() W365Ref {
	return n.Id
}

type Teacher struct {
	Id               W365Ref
	Type             string
	Name             string
	Shortcut         string
	Firstname        string
	Absences         []db.TimeSlot
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxDays          int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
	LunchBreak       bool
}

func (n *Teacher) IdStr() W365Ref {
	return n.Id
}

type Subject struct {
	Id       W365Ref
	Type     string
	Name     string
	Shortcut string
}

func (n *Subject) IdStr() W365Ref {
	return n.Id
}

type Room struct {
	Id       W365Ref
	Type     string
	Name     string
	Shortcut string
	Absences []db.TimeSlot
}

func (n *Room) IdStr() W365Ref {
	return n.Id
}

type RoomGroup struct {
	Id       W365Ref
	Type     string
	Name     string
	Shortcut string
	Rooms    []W365Ref
}

func (n *RoomGroup) IdStr() W365Ref {
	return n.Id
}

type Class struct {
	Id               W365Ref
	Type             string
	Name             string
	Shortcut         string
	Level            int
	Letter           string
	Absences         []db.TimeSlot
	Divisions        []Division
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
	LunchBreak       bool
	ForceFirstHour   bool
}

func (n *Class) IdStr() W365Ref {
	return n.Id
}

//func (n *Class) Tag() string {
//	return fmt.Sprintf("%d%s", n.Level, n.Letter)
//}

type Group struct {
	Id       W365Ref
	Type     string
	Shortcut string
}

func (n *Group) IdStr() W365Ref {
	return n.Id
}

type Division struct {
	Name   string
	Groups []W365Ref
}

type Course struct {
	Id             W365Ref
	Type           string
	Subjects       []W365Ref // if present, will be converted to Subject
	Subject        W365Ref
	Groups         []W365Ref
	Teachers       []W365Ref
	PreferredRooms []W365Ref
}

func (n *Course) IdStr() W365Ref {
	return n.Id
}

type SuperCourse struct {
	Id      W365Ref
	Type    string
	Subject W365Ref
}

func (n *SuperCourse) IdStr() W365Ref {
	return n.Id
}

type SubCourse struct {
	Id             W365Ref
	Type           string
	SuperCourse    W365Ref
	Subjects       []W365Ref // if present, will be converted to Subject
	Subject        W365Ref
	Groups         []W365Ref
	Teachers       []W365Ref
	PreferredRooms []W365Ref
}

func (n *SubCourse) IdStr() W365Ref {
	return n.Id
}

type Lesson struct {
	Id         W365Ref
	Type       string
	Course     W365Ref
	Duration   int
	Day        int
	Hour       int
	Fixed      bool
	LocalRooms []W365Ref
}

func (n *Lesson) IdStr() W365Ref {
	return n.Id
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
