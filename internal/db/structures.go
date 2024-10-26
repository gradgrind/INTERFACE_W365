package db

// The structures used for the "database"
//TODO: Currently dealing only with the elements needed for the timetable

type DbRef int // Element reference

type Info struct {
	Institution        string
	FirstAfternoonHour int
	MiddayBreak        []int
	Reference          interface{}
}

type Day struct {
	Id        DbRef
	Name      string
	Tag       string
	Reference interface{}
}

type Hour struct {
	Id        DbRef
	Name      string
	Tag       string
	Start     string
	End       string
	Reference interface{}
}

type TimeSlot struct {
	Day  int
	Hour int
}

type Teacher struct {
	Id               DbRef
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
	Reference        interface{}
}

type Subject struct {
	Id        DbRef
	Name      string
	Tag       string
	Reference interface{}
}

type Room struct {
	Id           DbRef
	Name         string
	Tag          string
	NotAvailable []TimeSlot
	Reference    interface{}
}

type RoomGroup struct {
	Id        DbRef
	Name      string
	Tag       string
	Rooms     []DbRef
	Reference interface{}
}

type RoomChoiceGroup struct {
	Id        DbRef
	Name      string
	Tag       string
	Rooms     []DbRef
	Reference interface{}
}

type Class struct {
	Id               DbRef
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
	Reference        interface{}
}

type Group struct {
	Id        DbRef
	Tag       string
	Reference interface{}
}

type Division struct {
	Name      string
	Groups    []DbRef
	Reference interface{}
}

type Course struct {
	Id        DbRef
	Subject   DbRef
	Groups    []DbRef
	Teachers  []DbRef
	Room      DbRef // Room, RoomGroup or RoomChoiceGroup Element
	Reference interface{}
}

type SuperCourse struct {
	Id        DbRef
	Subject   DbRef
	Reference interface{}
}

type SubCourse struct {
	Id          DbRef
	SuperCourse DbRef
	Subject     DbRef
	Groups      []DbRef
	Teachers    []DbRef
	Room        DbRef // Room, RoomGroup or RoomChoiceGroup Element
	Reference   interface{}
}

type Lesson struct {
	Id        DbRef
	Course    DbRef
	Duration  int
	Day       int
	Hour      int
	Fixed     bool
	Rooms     []DbRef // only Room Elements
	Reference interface{}
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
