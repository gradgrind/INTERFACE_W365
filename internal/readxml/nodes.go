package readxml

import (
	"encoding/xml"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/w365tt"
)

// The structures used for reading a timetable-source file exported by W365.

type W365RefList string // "List" of Element references

type TTNode interface {
	IdStr() w365tt.W365Ref
}

type SchoolState struct {
	ActiveScenario w365tt.W365Ref `xml:",attr"`
	SchoolName     string         `xml:",attr"`
}

type Day struct {
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"`
	Name         string         `xml:",attr"`
	Shortcut     string         `xml:",attr"`
}

func (n *Day) IdStr() w365tt.W365Ref {
	return n.Id
}

type Hour struct {
	XMLName            xml.Name       `xml:"TimedObject"`
	Id                 w365tt.W365Ref `xml:",attr"`
	ListPosition       float32        `xml:",attr"`
	Name               string         `xml:",attr"`
	Shortcut           string         `xml:",attr"`
	Start              string         `xml:",attr"`
	End                string         `xml:",attr"`
	FirstAfternoonHour bool           `xml:",attr"`
	MiddayBreak        bool           `xml:",attr"`
}

func (n *Hour) IdStr() w365tt.W365Ref {
	return n.Id
}

type Absence struct {
	Id   w365tt.W365Ref `xml:",attr"`
	Day  int            `xml:"day,attr" `
	Hour int            `xml:"hour,attr"`
}

func (n *Absence) IdStr() w365tt.W365Ref {
	return n.Id
}

type Teacher struct {
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"`
	Name         string         `xml:",attr"`
	Shortcut     string         `xml:",attr"`
	Firstname    string         `xml:",attr"`
	Absences     W365RefList    `xml:",attr"`
	Categories   W365RefList    `xml:",attr"`
	//+	Color string  `xml:",attr"` // "#ffcc00"
	//+	Gender int `xml:",attr"`
	MinLessonsPerDay int `xml:",attr"`
	MaxLessonsPerDay int `xml:",attr"`
	MaxDays          int `xml:",attr"`
	MaxGapsPerDay    int `xml:"MaxWindowsPerDay,attr"`
	//TODO: I have found MaxGapsPerWeek more useful
	MaxAfternoons int `xml:"NumberOfAfterNoonDays,attr"`
}

func (n *Teacher) IdStr() w365tt.W365Ref {
	return n.Id
}

type Subject struct {
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"`
	Name         string         `xml:",attr"`
	Shortcut     string         `xml:",attr"`
	//+ Color string  `xml:",attr"` // "#ffcc00"
}

func (n *Subject) IdStr() w365tt.W365Ref {
	return n.Id
}

type Room struct {
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"`
	Name         string         `xml:",attr"`
	Shortcut     string         `xml:",attr"`
	// The "Shortcut" can be very long when RoomGroups is not empty.
	// Name seems to be empty in these cases.
	Absences   W365RefList `xml:",attr"`
	Categories W365RefList `xml:",attr"`
	RoomGroups W365RefList `xml:"RoomGroup,attr"`
	// When RoomGroups is not empty, the "Room" is a room-group. In this
	// case ListPosition seems to be set to -1.
	//+ Capacity int `xml:"capacity,attr"`
	//+ Color string  `xml:",attr"` // "#ffcc00"
}

func (n *Room) IdStr() w365tt.W365Ref {
	return n.Id
}

type Class struct {
	XMLName      xml.Name       `xml:"Grade"`
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"`
	Name         string         `xml:",attr"` // Unused?
	//Shortcut         string   `xml:",attr"` // Unused?
	Level            int         `xml:",attr"`
	Letter           string      `xml:",attr"`
	Absences         W365RefList `xml:",attr"`
	Categories       W365RefList `xml:",attr"`
	ForceFirstHour   bool        `xml:",attr"`
	Divisions        W365RefList `xml:"GradePartitions,attr"`
	Groups           W365RefList `xml:",attr"` // Superfluous?
	MinLessonsPerDay int         `xml:",attr"`
	MaxLessonsPerDay int         `xml:",attr"`
	MaxAfternoons    int         `xml:"NumberOfAfterNoonDays,attr"`
	//+ ClassTeachers string `xml:"ClassTeacher,attr"`
	//+ Color string  `xml:",attr"` // "#ffcc00"
	//TODO: Implement in W365?
	//+ MaxGapsPerWeek    int `xml:"MaxWindowsPerWeek,attr"`
}

func (n *Class) IdStr() w365tt.W365Ref {
	return n.Id
}

func (n *Class) Tag() string {
	return fmt.Sprintf("%d%s", n.Level, n.Letter)
}

type Group struct {
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"` // Is this used?
	Name         string         `xml:",attr"` // How is this used?
	Shortcut     string         `xml:",attr"` // Presumably the primary tag ("eindeutig")
	//+ NumberOfStudents int     `xml:",attr"` // Is this used?
	//+ Color string  `xml:",attr"` // "#ffcc00" // Is this used?
}

func (n *Group) IdStr() w365tt.W365Ref {
	return n.Id
}

type Division struct {
	XMLName      xml.Name       `xml:"GradePartiton"`
	Id           w365tt.W365Ref `xml:",attr"`
	ListPosition float32        `xml:",attr"` // Is this used?
	Name         string         `xml:",attr"`
	Shortcut     string         `xml:",attr"` // Is this used?
	Groups       W365RefList    `xml:",attr"`
}

func (n *Division) IdStr() w365tt.W365Ref {
	return n.Id
}

type Course struct {
	Id                w365tt.W365Ref `xml:",attr"`
	ListPosition      float32        `xml:",attr"` // Is this used?
	Name              string         `xml:",attr"` // Is this used?
	Shortcut          string         `xml:",attr"` // Is this used?
	Subjects          W365RefList    `xml:",attr"` // can be more than one!
	Groups            W365RefList    `xml:",attr"` // either a Group or a Class?
	Teachers          W365RefList    `xml:",attr"`
	DoubleLessonMode  string         `xml:",attr"` // one course has "2,3"!
	HoursPerWeek      float32        `xml:",attr"`
	SplitHoursPerWeek string         `xml:",attr"` // "", "2+2+2+2+2" or "2+"
	PreferredRooms    W365RefList    `xml:",attr"`
	Categories        W365RefList    `xml:",attr"` // Is this used?
	EpochWeeks        float32        `xml:",attr"` // Is this relevant?
	//+ Color string  `xml:",attr"` // "#ffcc00" // Is this used?

	// These seem to be empty always. Are they relevant?
	//+ EpochSlots        string `xml:",attr"` // What is this?
	//+ SplitEpochWeeks   string `xml:",attr"` // What is this?
}

func (n *Course) IdStr() w365tt.W365Ref {
	return n.Id
}

type EpochPlanCourse struct {
	Id               w365tt.W365Ref `xml:",attr"`
	ListPosition     float32        `xml:",attr"` // Is this used?
	Name             string         `xml:",attr"` // Is this used?
	Shortcut         string         `xml:",attr"` // Is this used?
	Subjects         W365RefList    `xml:",attr"` // can be more than one!
	Groups           W365RefList    `xml:",attr"` // either a Group or a Class?
	Teachers         W365RefList    `xml:",attr"`
	DoubleLessonMode string         `xml:",attr"` // often "1,2"
	PreferredRooms   W365RefList    `xml:",attr"`
	Categories       W365RefList    `xml:",attr"` // Is this used?
	//+ Color string  `xml:",attr"` // "#ffcc00" // Is this used?

	//+ HoursPerWeek     float32 `xml:",attr"` // always 0.0?
	//+ EpochWeeks float32 `xml:",attr"` // always 0.0?

	// These seem to be always empty. Are they relevant?
	//+ EpochSlots        string `xml:",attr"` // What is this?
	//+ SplitEpochWeeks   string `xml:",attr"` // What is this?
	//+ SplitHoursPerWeek string `xml:",attr"` // What is this?
}

func (n *EpochPlanCourse) IdStr() w365tt.W365Ref {
	return n.Id
}

type Lesson struct {
	Id           w365tt.W365Ref `xml:",attr"`
	Course       w365tt.W365Ref `xml:",attr"`
	Day          int            `xml:",attr"`
	Hour         int            `xml:",attr"`
	DoubleLesson bool           `xml:",attr"` // What exactly does this mean?
	Fixed        bool           `xml:",attr"`
	Fractions    W365RefList    `xml:",attr"`
	LocalRooms   W365RefList    `xml:",attr"`
	EpochPlan    w365tt.W365Ref `xml:",attr"` // What is this? Not relevant?
	// If this entry is not empty, the Course field may be an EpochPlanCourse ... or nothing!
	EpochPlanGrade w365tt.W365Ref `xml:",attr"` // What is this?
}

func (n *Lesson) IdStr() w365tt.W365Ref {
	return n.Id
}

type Fraction struct {
	Id          w365tt.W365Ref `xml:",attr"`
	SuperGroups W365RefList    `xml:",attr"`
}

func (n *Fraction) IdStr() w365tt.W365Ref {
	return n.Id
}

type W365TTXML struct {
	XMLName          xml.Name `xml:"File"`
	Path             string
	SchoolState      SchoolState
	Days             []Day             `xml:"Day"`
	Hours            []Hour            `xml:"TimedObject"`
	Absences         []Absence         `xml:"Absence"`
	Teachers         []Teacher         `xml:"Teacher"`
	Subjects         []Subject         `xml:"Subject"`
	Rooms            []Room            `xml:"Room"`
	Classes          []Class           `xml:"Grade"`
	Groups           []Group           `xml:"Group"`
	Divisions        []Division        `xml:"GradePartiton"`
	Courses          []Course          `xml:"Course"`
	EpochPlanCourses []EpochPlanCourse `xml:"EpochPlanCourse"`
	Lessons          []Lesson          `xml:"Lesson"`
	Fractions        []Fraction        `xml:"Fraction"`
}
