use serde::{Deserialize, Serialize};
use serde_json::Value;

// The structures used as a general database.

pub type DbRef = u32; // Element reference

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct TimeSlot {
    pub Day: i32,
    pub Hour: i32
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Info {
	pub Institution:        String,
	pub FirstAfternoonHour: i32,
	pub MiddayBreak:        Vec<i32>,
    pub Reference:          Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Day {
	pub Id:         DbRef,
	pub Name:       String,
	pub Tag:        String,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Hour {
	pub Id:         DbRef,
	pub Name:       String,
	pub Tag:        String,
	pub Start:      String,
	pub End:        String,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Teacher {
	pub Id:                 DbRef,
	pub Name:               String,
	pub Tag:                String,
	pub Firstname:          String,
	pub NotAvailable:       Vec<TimeSlot>,
	pub MinLessonsPerDay:   i32,
	pub MaxLessonsPerDay:   i32,
	pub MaxDays:            i32,
	pub MaxGapsPerDay:      i32,
	pub MaxGapsPerWeek:     i32,
	pub MaxAfternoons:      i32,
	pub LunchBreak:         bool,
    pub Reference:          Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Subject {
	pub Id:         DbRef,
	pub Name:       String,
	pub Tag:        String,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Room {
	pub Id:             DbRef,
	pub Name:           String,
	pub Tag:            String,
	pub NotAvailable:   Vec<TimeSlot>,
    pub Reference:      Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct RoomGroup {
	pub Id:         DbRef,
	pub Name:       String,
	pub Tag:        String,
	pub Rooms:      Vec<DbRef>,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct RoomChoiceGroup {
	pub Id:         DbRef,
	pub Name:       String,
	pub Tag:        String,
	pub Rooms:      Vec<DbRef>,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Class {
	pub Id:                 DbRef,
	pub Name:               String,
	pub Tag:                String,
	pub Level:              i32,
	pub Letter:             String,
	pub Absences:           Vec<TimeSlot>,
	pub Divisions:          Vec<Division>,
	pub MinLessonsPerDay:   i32,
	pub MaxLessonsPerDay:   i32,
	pub MaxGapsPerDay:      i32,
	pub MaxGapsPerWeek:     i32,
	pub MaxAfternoons:      i32,
	pub LunchBreak:         bool,
	pub ForceFirstHour:     bool,
    pub Reference:          Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Group {
	pub Id:         DbRef,
	pub Tag:        String,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Division {
	pub Name:       String,
	pub Groups:     Vec<DbRef>,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Course {
	pub Id:         DbRef,
	pub Subject:    DbRef,
	pub Groups:     Vec<DbRef>,
	pub Teachers:   Vec<DbRef>,
	pub Room:       DbRef, // Room, RoomGroup or RoomChoiceGroup Element
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct SuperCourse {
	pub Id:         DbRef,
	pub Subject:    DbRef,
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct SubCourse {
	pub Id:             DbRef,
	pub SuperCourse:    DbRef,
	pub Subject:        DbRef,
	pub Groups:         Vec<DbRef>,
	pub Teachers:       Vec<DbRef>,
	pub Room:           DbRef, // Room, RoomGroup or RoomChoiceGroup Element
    pub Reference:      Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct Lesson {
	pub Id:         DbRef,
	pub Course:     DbRef, // Course or SuperCourse Element
	pub Duration:   i32,
	pub Day:        i32,
	pub Hour:       i32,
	pub Fixed:      bool,
	pub Rooms:      Vec<DbRef>, // only Room Elements
    pub Reference:  Value
}

#[allow(nonstandard_style)]
#[derive(Serialize, Deserialize, Debug)]
pub struct DbTopLevel {
	pub Info:               Info,
	pub Days:               Vec<Day>,
	pub Hours:              Vec<Hour>,
	pub Teachers:           Vec<Teacher>,
	pub Subjects:           Vec<Subject>,
	pub Rooms:              Vec<Room>,
    pub RoomGroups:         Vec<RoomGroup>,
    pub RoomChoiceGroups:   Vec<RoomGroup>,
	pub Classes:            Vec<Class>,
	pub Groups:             Vec<Group>,
	pub Courses:            Vec<Course>,
	pub SuperCourses:       Vec<SuperCourse>,
	pub SubCourses:         Vec<SubCourse>,
	pub Lessons:            Vec<Lesson>,
	pub Constraints:        Value
}
