use serde::{Deserialize, Serialize};
use serde_json::Result;

// The structures used for reading a timetable-source file exported by W365.

type W365Ref = String; // Element reference

struct Info {
	SchoolName:         String,
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
	Name:               String,
	Shortcut:           String,
	Start:              String,
	End:                String,
	FirstAfternoonHour: bool,   // default = false
	MiddayBreak:        bool    // default = false
}

fn default_m1() -> i32 { -1 }


//TODO: Use a tuple for a timeslot?

#[derive(Deserialize, Serialize)]
struct Teacher {
	Id:               W365Ref,
	Name:             String,
	Shortcut:         String,
	Firstname:        String,
	Absences:         Vec<(usize, usize)>,
	#[serde(default = "default_m1")]
    MinLessonsPerDay: i32,
	#[serde(default = "default_m1")]
    MaxLessonsPerDay: i32,
	#[serde(default = "default_m1")]
    MaxDays:          i32,
	#[serde(default = "default_m1")]
    MaxGapsPerDay:    i32,
	#[serde(default = "default_m1")]
    MaxGapsPerWeek:   i32,
	#[serde(default = "default_m1")]
	MaxAfternoons:    i32,
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
	Absences:   Vec<(usize, usize)>
}

struct RoomGroup {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String,
	Rooms:      Vec<W365Ref>
}

#[derive(Deserialize, Serialize)]
struct Class {
	Id:                 W365Ref,
	Name:               String,
	Shortcut:           String,
	Level:              i32,
	Letter:             String,
	Absences:           Vec<(usize, usize)>,
	Divisions:          Vec<Division>,
	#[serde(default = "default_m1")]
    MinLessonsPerDay:   i32,
	#[serde(default = "default_m1")]
    MaxLessonsPerDay:   i32,
	#[serde(default = "default_m1")]
    MaxGapsPerDay:      i32,
	#[serde(default = "default_m1")]
    MaxGapsPerWeek:     i32,
	#[serde(default = "default_m1")]
	MaxAfternoons:      i32,
	LunchBreak:         bool,
	ForceFirstHour:     bool
}

struct Group {
	Id:         W365Ref,
	Shortcut:   String
}

#[derive(Deserialize, Serialize)]
struct Division {
	Id:         W365Ref, // even though it is not top-level
	Name:       String,
	Groups:     Vec<W365Ref>
}

struct Course {
	Id:             W365Ref,
	Subjects:       Vec<W365Ref>,
	Subject:        W365Ref,
	Groups:         Vec<W365Ref>,
	Teachers:       Vec<W365Ref>,
	PreferredRooms: Vec<W365Ref>
}

struct SuperCourse {
	Id:         W365Ref,
	Subject:    W365Ref
}

#[derive(Deserialize, Serialize)]
struct SubCourse {
	Id:             W365Ref,
	SuperCourse:    W365Ref,
	#[serde(default)]
    Subjects:       Vec<W365Ref>,
	#[serde(default)]
    Subject:        W365Ref,
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
