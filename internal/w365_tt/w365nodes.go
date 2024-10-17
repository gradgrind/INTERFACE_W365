package w365_tt

import "encoding/xml"

// The structures used for reading a timetable-source file exported by W365.

type Day struct {
	Id           string  `xml:",attr"`
	ListPosition float32 `xml:",attr"`
	Name         string  `xml:",attr"`
	Shortcut     string  `xml:",attr"`
}

type Hour struct {
	XMLName            xml.Name `xml:"TimedObject"`
	Id                 string   `xml:",attr"`
	ListPosition       float32  `xml:",attr"`
	Name               string   `xml:",attr"`
	Shortcut           string   `xml:",attr"`
	Start              string   `xml:",attr"`
	End                string   `xml:",attr"`
	FirstAfternoonHour bool     `xml:",attr"`
	MiddayBreak        bool     `xml:",attr"`
}

type Absence struct {
	Id   string `xml:",attr"`
	Day  int    `xml:"day,attr" `
	Hour int    `xml:"hour,attr"`
}

type Teacher struct {
	Id           string  `xml:",attr"`
	ListPosition float32 `xml:",attr"`
	Name         string  `xml:",attr"`
	Shortcut     string  `xml:",attr"`
	Firstname    string  `xml:",attr"`
	Absences     string  `xml:",attr"`
	Categories   string  `xml:",attr"`
	//+	Color string  `xml:",attr"` // "#ffcc00"
	//+	Gender int `xml:",attr"`
	MinLessonsPerDay int `xml:",attr"`
	MaxLessonsPerDay int `xml:",attr"`
	MaxDays          int `xml:",attr"`
	MaxGapsPerDay    int `xml:"MaxWindowsPerDay,attr"`
	//TODO: I have found MaxGapsPerWeek more useful
	MaxAfternoons int `xml:"NumberOfAfterNoonDays,attr"`
}

type Subject struct {
	Id           string  `xml:",attr"`
	ListPosition float32 `xml:",attr"`
	Name         string  `xml:",attr"`
	Shortcut     string  `xml:",attr"`
	//+ Color string  `xml:",attr"` // "#ffcc00"
}

type Room struct {
	Id           string  `xml:",attr"`
	ListPosition float32 `xml:",attr"`
	Name         string  `xml:",attr"`
	Shortcut     string  `xml:",attr"`
	// The "Shortcut" can be very long when RoomGroups is not empty.
	// Name seems to be empty in these cases.
	Absences   string `xml:",attr"`
	Categories string `xml:",attr"`
	RoomGroups string `xml:"RoomGroup,attr"`
	// When RoomGroups is not empty, the "Room" is a room-group. In this
	// case ListPosition seems to be set to -1.
	//+ Capacity int `xml:"capacity,attr"`
	//+ Color string  `xml:",attr"` // "#ffcc00"
}

type Class struct {
	XMLName          xml.Name `xml:"Grade"`
	Id               string   `xml:",attr"`
	ListPosition     float32  `xml:",attr"`
	Name             string   `xml:",attr"` // How is this used?
	Shortcut         string   `xml:",attr"` // Presumably the primary tag ("eindeutig")
	Level            int      `xml:",attr"`
	Letter           string   `xml:",attr"`
	Absences         string   `xml:",attr"`
	Categories       string   `xml:",attr"`
	ForceFirstHour   bool     `xml:",attr"`
	Divisions        string   `xml:"GradePartitions,attr"`
	Groups           string   `xml:",attr"`
	MinLessonsPerDay int      `xml:",attr"`
	MaxLessonsPerDay int      `xml:",attr"`
	MaxAfternoons    int      `xml:"NumberOfAfterNoonDays,attr"`
	//+ ClassTeachers string `xml:"ClassTeacher,attr"`
	//+ Color string  `xml:",attr"` // "#ffcc00"
	//TODO: Implement in W365?
	//+ MaxGapsPerWeek    int `xml:"MaxWindowsPerWeek,attr"`
}

type Group struct {
	Id           string  `xml:",attr"`
	ListPosition float32 `xml:",attr"` // Is this used?
	Name         string  `xml:",attr"` // How is this used?
	Shortcut     string  `xml:",attr"` // Presumably the primary tag ("eindeutig")
	Absences     string  `xml:",attr"` // Is this used?
	Categories   string  `xml:",attr"` // Is this used?
	//+ NumberOfStudents int     `xml:",attr"` // Is this used?
	//+ Color string  `xml:",attr"` // "#ffcc00" // Is this used?
}

type Division struct {
	XMLName      xml.Name `xml:"GradePartiton"`
	Id           string   `xml:",attr"`
	ListPosition float32  `xml:",attr"` // Is this used?
	Name         string   `xml:",attr"`
	Shortcut     string   `xml:",attr"` // Is this used?
	Groups       string   `xml:",attr"`
}

type Course struct {
	Id                string  `xml:",attr"`
	ListPosition      float32 `xml:",attr"` // Is this used?
	Name              string  `xml:",attr"` // Is this used?
	Shortcut          string  `xml:",attr"` // Is this used?
	Subjects          string  `xml:",attr"` // can be more than one!
	Groups            string  `xml:",attr"` // either a Group or a Class?
	Teachers          string  `xml:",attr"`
	DoubleLessonMode  string  `xml:",attr"` // one course has "2,3"!
	HoursPerWeek      float32 `xml:",attr"`
	SplitHoursPerWeek string  `xml:",attr"` // "", "2+2+2+2+2" or "2+"
	PreferredRooms    string  `xml:",attr"`
	Categories        string  `xml:",attr"` // Is this used?
	EpochWeeks        float32 `xml:",attr"` // Is this relevant?
	//+ Color string  `xml:",attr"` // "#ffcc00" // Is this used?

	// These seem to be empty always. Are they relevant?
	//+ EpochSlots        string `xml:",attr"` // What is this?
	//+ SplitEpochWeeks   string `xml:",attr"` // What is this?
}

type EpochPlanCourse struct {
	Id               string  `xml:",attr"`
	ListPosition     float32 `xml:",attr"` // Is this used?
	Name             string  `xml:",attr"` // Is this used?
	Shortcut         string  `xml:",attr"` // Is this used?
	Subjects         string  `xml:",attr"` // can be more than one!
	Groups           string  `xml:",attr"` // either a Group or a Class?
	Teachers         string  `xml:",attr"`
	DoubleLessonMode string  `xml:",attr"` // often "1,2"
	PreferredRooms   string  `xml:",attr"`
	Categories       string  `xml:",attr"` // Is this used?
	//+ Color string  `xml:",attr"` // "#ffcc00" // Is this used?

	//+ HoursPerWeek     float32 `xml:",attr"` // always 0.0?
	//+ EpochWeeks float32 `xml:",attr"` // always 0.0?

	// These seem to be always empty. Are they relevant?
	//+ EpochSlots        string `xml:",attr"` // What is this?
	//+ SplitEpochWeeks   string `xml:",attr"` // What is this?
	//+ SplitHoursPerWeek string `xml:",attr"` // What is this?
}

type Lesson struct {
	Id           string `xml:",attr"`
	Course       string `xml:",attr"`
	Day          int    `xml:",attr"`
	Hour         int    `xml:",attr"`
	DoubleLesson bool   `xml:",attr"` // What exactly does this mean?
	Fixed        bool   `xml:",attr"`
	Fractions    string `xml:",attr"`
	LocalRooms   string `xml:",attr"`
	EpochPlan    string `xml:",attr"` // What is this? Not relevant?
	// If this entry is not empty, the Course field may be an EpochPlanCourse ... or nothing!
	EpochPlanGrade string `xml:",attr"` // What is this?
}

type Fraction struct {
	Id          string `xml:",attr"`
	SuperGroups string `xml:",attr"`
}

type W365TT struct {
	XMLName          xml.Name          `xml:"File"`
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
