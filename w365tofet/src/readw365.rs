use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::fs;
use std::error::Error;

pub fn read_w365(jsonpath: &String) -> Result<W365TopLevel,Box<dyn Error>> {

    let w365json = fs::read_to_string(jsonpath)?;
//        .expect("Couldn't read input file");

    let w365data: W365TopLevel = serde_json::from_str(&w365json)?;
//        .expect("Couldn't parse JSON");
    Ok(w365data)
}

// The structures used for reading a timetable-source file exported by W365.

pub type W365Ref = String; // Element reference

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct TimeSlot {
    pub Day: i32,
    pub Hour: i32
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Info {
	pub SchoolName:         String,
	pub Scenario:           W365Ref,
	pub FirstAfternoonHour: i32,
	pub MiddayBreak:        Vec<i32>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Day {
	pub Id:         W365Ref,
	pub Name:       String,
	pub Shortcut:   String
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Hour {
	pub Id:                 W365Ref,
	pub Name:               String,
	pub Shortcut:           String,
	pub Start:              String,
	pub End:                String,
	pub FirstAfternoonHour: bool,   // default = false
	pub MiddayBreak:        bool    // default = false
}

fn default_m1() -> i32 { -1 }

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Teacher {
	pub Id:               W365Ref,
	pub Name:             String,
	pub Shortcut:         String,
	pub Firstname:        String,
	pub Absences:         Vec<TimeSlot>,
	#[serde(default = "default_m1")]
    pub MinLessonsPerDay: i32,
	#[serde(default = "default_m1")]
    pub MaxLessonsPerDay: i32,
	#[serde(default = "default_m1")]
    pub MaxDays:          i32,
	#[serde(default = "default_m1")]
    pub MaxGapsPerDay:    i32,
	#[serde(default = "default_m1")]
    pub MaxGapsPerWeek:   i32,
	#[serde(default = "default_m1")]
	pub MaxAfternoons:    i32,
	pub LunchBreak:       bool
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Subject {
	pub Id:         W365Ref,
	pub Name:       String,
	pub Shortcut:   String
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Room {
	pub Id:         W365Ref,
	pub Name:       String,
	pub Shortcut:   String,
	pub Absences:   Vec<TimeSlot>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct RoomGroup {
	pub Id:         W365Ref,
	pub Name:       String,
	pub Shortcut:   String,
	pub Rooms:      Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Class {
	pub Id:                 W365Ref,
	pub Name:               String,
	pub Shortcut:           String,
	pub Level:              i32,
	pub Letter:             String,
	pub Absences:           Vec<TimeSlot>,
	pub Divisions:          Vec<Division>,
	#[serde(default = "default_m1")]
    pub MinLessonsPerDay:   i32,
	#[serde(default = "default_m1")]
    pub MaxLessonsPerDay:   i32,
	#[serde(default = "default_m1")]
    pub MaxGapsPerDay:      i32,
	#[serde(default = "default_m1")]
    pub MaxGapsPerWeek:     i32,
	#[serde(default = "default_m1")]
	pub MaxAfternoons:      i32,
	pub LunchBreak:         bool,
	pub ForceFirstHour:     bool
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Group {
	pub Id:         W365Ref,
	pub Shortcut:   String
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Division {
	pub Id:         W365Ref, // even though it is not top-level
	pub Name:       String,
	pub Groups:     Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Course {
	pub Id:             W365Ref,
	#[serde(default)]
	pub Subjects:       Vec<W365Ref>,
	#[serde(default)]
    pub Subject:        W365Ref,
	pub Groups:         Vec<W365Ref>,
	pub Teachers:       Vec<W365Ref>,
	pub PreferredRooms: Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct SuperCourse {
	pub Id:         W365Ref,
	pub Subject:    W365Ref
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct SubCourse {
	pub Id:             W365Ref,
	pub SuperCourse:    W365Ref,
	#[serde(default)]
    pub Subjects:       Vec<W365Ref>,
	#[serde(default)]
    pub Subject:        W365Ref,
	pub Groups:         Vec<W365Ref>,
	pub Teachers:       Vec<W365Ref>,
	pub PreferredRooms: Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Lesson {
	pub Id:         W365Ref,
	pub Course:     W365Ref,
	pub Duration:   i32,
	pub Day:        i32,
	pub Hour:       i32,
	pub Fixed:      bool,
	pub LocalRooms: Vec<W365Ref>
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct W365TopLevel {
	pub W365TT:         Info,
	pub Days:           Vec<Day>,
	pub Hours:          Vec<Hour>,
	pub Teachers:       Vec<Teacher>,
	pub Subjects:       Vec<Subject>,
	pub Rooms:          Vec<Room>,
    pub RoomGroups:     Vec<RoomGroup>,
	pub Classes:        Vec<Class>,
	pub Groups:         Vec<Group>,
	pub Courses:        Vec<Course>,
	pub SuperCourses:   Vec<SuperCourse>,
	pub SubCourses:     Vec<SubCourse>,
	pub Lessons:        Vec<Lesson>,
	pub Constraints:    Value
}
