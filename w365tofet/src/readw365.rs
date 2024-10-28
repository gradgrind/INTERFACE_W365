use serde::{Deserialize, Serialize};
use serde_json::{Result, Value};
use std::fs;

pub fn read_w365(jsonpath: String) {

    let w365json = fs::read_to_string(jsonpath)
        .expect("Couldn't read input file");

    let w365data: W365TopLevel = serde_json::from_str(&w365json)
        .expect("Couldn't parse JSON");
    println!("{:#?}", w365data);
}
// The structures used for reading a timetable-source file exported by W365.

type W365Ref = String; // Element reference

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct TimeSlot {
    Day: i32,
    Hour: i32
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Info {
	SchoolName:         String,
	Scenario:           W365Ref,
	FirstAfternoonHour: i32,
	MiddayBreak:        Vec<i32>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Day {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
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

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Teacher {
	Id:               W365Ref,
	Name:             String,
	Shortcut:         String,
	Firstname:        String,
	Absences:         Vec<TimeSlot>,
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

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Subject {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Room {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String,
	Absences:   Vec<TimeSlot>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct RoomGroup {
	Id:         W365Ref,
	Name:       String,
	Shortcut:   String,
	Rooms:      Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Class {
	Id:                 W365Ref,
	Name:               String,
	Shortcut:           String,
	Level:              i32,
	Letter:             String,
	Absences:           Vec<TimeSlot>,
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

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Group {
	Id:         W365Ref,
	Shortcut:   String
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Division {
	Id:         W365Ref, // even though it is not top-level
	Name:       String,
	Groups:     Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Course {
	Id:             W365Ref,
	#[serde(default)]
	Subjects:       Vec<W365Ref>,
	#[serde(default)]
    Subject:        W365Ref,
	Groups:         Vec<W365Ref>,
	Teachers:       Vec<W365Ref>,
	PreferredRooms: Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct SuperCourse {
	Id:         W365Ref,
	Subject:    W365Ref
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
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

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct Lesson {
	Id:         W365Ref,
	Course:     W365Ref,
	Duration:   i32,
	Day:        i32,
	Hour:       i32,
	Fixed:      bool,
	LocalRooms: Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
struct W365TopLevel {
	W365TT:         Info,
	Days:           Vec<Day>,
	Hours:          Vec<Hour>,
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
	Constraints:    Value
}
