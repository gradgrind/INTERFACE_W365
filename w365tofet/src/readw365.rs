

// The structures used for reading a timetable-source file exported by W365.

type W365Ref = String; // Element reference

struct Info {
	SchoolName:         string,
	Scenario:           W365Ref,
	FirstAfternoonHour: i32,
	MiddayBreak:        Vec<i32>
}

struct Day {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String
}

struct Hour {
	Id:                 W365Ref,
	Name:               string,
	Shortcut:           string,
	Start:              string,
	End:                string,
	FirstAfternoonHour: bool,   // default = false
	MiddayBreak:        bool    // default = false
}

struct Teacher {
	Id:               W365Ref,
	Name:             String,
	Shortcut:         String,
	Firstname:        String,
	Absences:         []db.TimeSlot,
	MinLessonsPerDay: interface{} `json:",omitempty"`
	MaxLessonsPerDay: interface{} `json:",omitempty"`
	MaxDays:          interface{} `json:",omitempty"`
	MaxGapsPerDay:    interface{} `json:",omitempty"`
	MaxGapsPerWeek:   interface{} `json:",omitempty"`
	MaxAfternoons:    interface{} `json:",omitempty"`
	LunchBreak:       bool
}

struct Subject {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String
}

struct Room {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String,
	Absences:   []db.TimeSlot
}

struct RoomGroup {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String,
	Rooms:      Vec<W365Ref>
}

struct Class {
	Id:                 W365Ref,
	Name:               String,
	Shortcut:           String,
	Level:              i32,
	Letter:             String,
	Absences:           []db.TimeSlot
	Divisions:          []Division
	MinLessonsPerDay:   interface{} `json:",omitempty"`
	MaxLessonsPerDay:   interface{} `json:",omitempty"`
	MaxGapsPerDay:      interface{} `json:",omitempty"`
	MaxGapsPerWeek:     interface{} `json:",omitempty"`
	MaxAfternoons:      interface{} `json:",omitempty"`
	LunchBreak:         bool,
	ForceFirstHour:     bool
}

struct Group {
	Id:         W365Ref,
	Shortcut:   String
}

struct Division {
	Id:         W365Ref, // even though it is not top-level
	Name:       String,
	Groups:     Vec<W365Ref>
}

struct Course {
	Id:             W365Ref,
	Subjects:       Vec<W365Ref>, `json:",omitempty"`
	Subject:        W365Ref,      `json:",omitempty"`
	Groups:         Vec<W365Ref>,
	Teachers:       Vec<W365Ref>,
	PreferredRooms: Vec<W365Ref>
}

struct SuperCourse {
	Id:         W365Ref,
	Subject:    W365Ref
}

struct SubCourse {
	Id:             W365Ref,
	SuperCourse:    W365Ref,
	Subjects:       Vec<W365Ref>, `json:",omitempty"`
	Subject:        W365Ref,   `json:",omitempty"`
	Groups:         Vec<W365Ref>,
	Teachers:       Vec<W365Ref>,
	PreferredRooms: Vec<W365Ref>
}

struct Lesson {
	Id:         W365Ref,
	Course:     W365Ref,
	Duration:   i32,
	Day:        i32,
	Hour:       i32,
	Fixed:      bool,
	LocalRooms: Vec<W365Ref>
}

struct W365TopLevel {
	W365TT:         Info,
	Days:           Day,
	Hours:          Hour,
	Teachers:       Vec<Teacher>,
	Subjects:       Vec<Subject>,
	Rooms:          Vec<Room>,
	RoomGroups:     Vec<RoomGroup>,
	Classes:        Vec<Class>,
	Groups:         Vec<Group>,
	Courses:        Vec<Course>,
	SuperCourses:   Vec<SuperCourse>,
	SubCourses:     Vec<SubCourse>,
	Lessons:        Vec<Lesson>,
	Constraints:    map[string]interface{}
}
