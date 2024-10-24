package db

// The structures used for the "database"
//TODO: Currently dealing only with the elements needed for the timetable

const (
	TypeDAY         string = "Day"
	TypeHOUR        string = "Hour"
	TypeTEACHER     string = "Teacher"
	TypeSUBJECT     string = "Subject"
	TypeROOM        string = "Room"
	TypeROOMGROUP   string = "RoomGroup"
	TypeCLASS       string = "Class"
	TypeGROUP       string = "Group"
	TypeCOURSE      string = "Course"
	TypeSUPERCOURSE string = "SuperCourse"
	TypeSUBCOURSE   string = "SubCourse"
	TypeLESSON      string = "Lesson"
)

type DbRef int // Element reference

type TTNode interface {
	// An interface for Top-Level-Elements with Id field
	IdStr() DbRef
}

type Info struct {
	FirstAfternoonHour int
	MiddayBreak        []int
}

type Day struct {
	Id   DbRef
	Type string
	Name string
	Tag  string
}

func (n *Day) IdStr() DbRef {
	return n.Id
}

type Hour struct {
	Id    DbRef
	Type  string
	Name  string
	Tag   string
	Start string
	End   string
}

func (n *Hour) IdStr() DbRef {
	return n.Id
}

type TimeSlot struct {
	Day  int
	Hour int
}

type Teacher struct {
	Id               DbRef
	Type             string
	Name             string
	Tag              string
	Firstname        string
	NotAvailable     []TimeSlot
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxDays          int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
	LunchBreak       bool
}

func (n *Teacher) IdStr() DbRef {
	return n.Id
}

type Subject struct {
	Id   DbRef
	Type string
	Name string
	Tag  string
}

func (n *Subject) IdStr() DbRef {
	return n.Id
}

type Room struct {
	Id           DbRef
	Type         string
	Name         string
	Tag          string
	NotAvailable []TimeSlot
}

func (n *Room) IdStr() DbRef {
	return n.Id
}

type RoomChoiceGroup struct {
	Id    DbRef
	Type  string
	Name  string
	Tag   string
	Rooms []DbRef
}

func (n *RoomChoiceGroup) IdStr() DbRef {
	return n.Id
}

type Class struct {
	Id               DbRef
	Type             string
	Name             string
	Tag              string
	Level            int
	Letter           string
	NotAvailable     []TimeSlot
	Divisions        []Division
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
	LunchBreak       bool
	ForceFirstHour   bool
}

func (n *Class) IdStr() DbRef {
	return n.Id
}

//func (n *Class) Tag() string {
//	return fmt.Sprintf("%d%s", n.Level, n.Letter)
//}

type Group struct {
	Id   DbRef
	Type string
	Tag  string
}

func (n *Group) IdStr() DbRef {
	return n.Id
}

type Division struct {
	Name   string
	Groups []DbRef
}

type Course struct {
	Id      DbRef
	Type    string
	Subject DbRef
	Groups  []DbRef
	Teacher DbRef
	Rooms   []DbRef // Room and RoomChoiceGroup Elements permitted
}

func (n *Course) IdStr() DbRef {
	return n.Id
}

type SuperCourse struct {
	Id      DbRef
	Type    string
	Subject DbRef
}

func (n *SuperCourse) IdStr() DbRef {
	return n.Id
}

type SubCourse struct {
	Id          DbRef
	Type        string
	SuperCourse DbRef
	Subject     DbRef
	Groups      []DbRef
	Teachers    []DbRef
	Rooms       []DbRef // Room and RoomChoiceGroup Elements permitted
}

func (n *SubCourse) IdStr() DbRef {
	return n.Id
}

type Lesson struct {
	Id       DbRef
	Type     string
	Course   DbRef
	Duration int
	Day      int
	Hour     int
	Fixed    bool
	Rooms    []DbRef // only Room Elements
}

func (n *Lesson) IdStr() DbRef {
	return n.Id
}

type DbTopLevel struct {
	Info             Info
	Days             []Day
	Hours            []Hour
	Teachers         []Teacher
	Subjects         []Subject
	Rooms            []Room
	RoomChoiceGroups []RoomChoiceGroup
	Classes          []Class
	Groups           []Group
	Courses          []Course
	SuperCourses     []SuperCourse
	SubCourses       []SubCourse
	Lessons          []Lesson
	Constraints      map[string]interface{}
}
