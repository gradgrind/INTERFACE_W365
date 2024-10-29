use crate::readw365;
use crate::readw365::W365Ref;
use crate::db;
use crate::db::DbRef;
use std::collections::BTreeMap;
use serde_json::json;

type MapWDb = BTreeMap<W365Ref, DbRef>;
type MapWStr = BTreeMap<W365Ref, String>;
type MapStrDb = BTreeMap<String, DbRef>;
type MapDbStr = BTreeMap<DbRef, String>;


struct XData {
	//w365:         readw365::W365TopLevel,
	data:         db::DbTopLevel,
	dbi:          DbRef, // counter, for db indexes
	teachers:     MapWDb,
	subjects:     MapWDb,
	subjectmap:   MapWStr, // Subject Tag (Shortcut)
	rooms:        MapWDb,
	roomtag:      MapDbStr, // Room Tag (Shortcut)
	roomgroups:   MapWDb,
	roomchoices:  MapStrDb, // New RoomChoiceGroup name -> db Id
	pregroups:    MapWStr,
	groups:       MapWDb,
	classes:      MapWDb,
	courses:      MapWDb,
	subcourses:   MapWDb,
	supercourses: MapWDb,
	newsubjects:  MapStrDb // New Subject name -> db Id
}

pub fn w365_db(w365data: readw365::W365TopLevel)
        -> Result<db::DbTopLevel, String>
{
	let mut dbdata = XData{
        //w365:           w365data,
        data:           db::DbTopLevel::new(db::Info{
            Institution:        w365data.W365TT.SchoolName.clone(),
            FirstAfternoonHour: w365data.W365TT.FirstAfternoonHour,
            MiddayBreak:        w365data.W365TT.MiddayBreak.clone(),
            Reference:          json!(w365data.W365TT.Scenario)
        }),
        dbi:            0,
        teachers:       BTreeMap::new(),
        subjects:       BTreeMap::new(),
        subjectmap:     BTreeMap::new(),
        rooms:          BTreeMap::new(),
        roomtag:        BTreeMap::new(),
        roomgroups:     BTreeMap::new(),
        roomchoices:    BTreeMap::new(),
        pregroups:      BTreeMap::new(),
        groups:         BTreeMap::new(),
        classes:        BTreeMap::new(),
        courses:        BTreeMap::new(),
        subcourses:     BTreeMap::new(),
        supercourses:   BTreeMap::new(),
        newsubjects:    BTreeMap::new()
    };

	add_days(&mut dbdata, &w365data.Days);
	add_hours(&mut dbdata, &w365data.Hours);
	add_teachers(&mut dbdata, &w365data.Teachers);
    Ok(dbdata.data)
}

fn next_id(dbdata: &mut XData) -> DbRef {
    dbdata.dbi += 1;
    dbdata.dbi
}

fn add_days(dbdata: &mut XData, days: &Vec<readw365::Day>) {
	for d in days.iter() {
        let id = next_id(dbdata);
		dbdata.data.Days.push(db::Day{
			Id:        id,
			Tag:       d.Shortcut.clone(),
			Name:      d.Name.clone(),
			Reference: json!(d.Id)
		})
	}
}

fn add_hours(dbdata: &mut XData, hours: &Vec<readw365::Hour>) {
	let mdbok = dbdata.data.Info.MiddayBreak.is_empty();
    let mut i = 0;
	for d in hours.iter() {
		if d.FirstAfternoonHour {
			dbdata.data.Info.FirstAfternoonHour = i;
		}
		if d.MiddayBreak {
			if mdbok {
				dbdata.data.Info.MiddayBreak.push(i);
			} else {
				eprintln!("*ERROR* MiddayBreak set in Info AND Hours")
			}
		}
        i += 1;
		let id = next_id(dbdata);
		dbdata.data.Hours.push(db::Hour{
			Id:        id,
			Tag:       if d.Shortcut.is_empty() {
                format!("({})", i)
            } else {
                d.Shortcut.clone()
            },
			Name:      d.Name.clone(),
			Start:     d.Start.clone(),
			End:       d.End.clone(),
			Reference: json!(d.Id)
		})
	}
}

fn clone_absences(alist: &Vec<db::TimeSlot>) -> Vec<db::TimeSlot> {
    let mut a: Vec<db::TimeSlot> = Vec::new();
    if !alist.is_empty() {
        for t in alist.iter() {
            a.push(db::TimeSlot{Day: t.Day, Hour: t.Hour});
        }
    }
    a
}

fn add_teachers(dbdata: &mut XData, teachers: &Vec<readw365::Teacher>) {
	for d in teachers.iter() {
		let id = next_id(dbdata);
		dbdata.data.Teachers.push(db::Teacher{
			Id:                 id,
            Tag:                d.Shortcut.clone(),
			Name:               d.Name.clone(),
			Firstname:          d.Firstname.clone(),
			NotAvailable:       clone_absences(&d.Absences),
			MinLessonsPerDay:   d.MinLessonsPerDay,
			MaxLessonsPerDay:   d.MaxLessonsPerDay,
			MaxDays:            d.MaxDays,
			MaxGapsPerDay:      d.MaxGapsPerDay,
			MaxGapsPerWeek:     d.MaxGapsPerWeek,
			MaxAfternoons:      d.MaxAfternoons,
			LunchBreak:         d.LunchBreak,
			Reference:          json!(d.Id)
		});
		dbdata.teachers.insert(d.Id.clone(), id);
	}
}


/*

func LoadJSON(jsonpath string) db.DbTopLevel {
	dbdata := xData{
		w365: ReadJSON(jsonpath),
		data: db.DbTopLevel{},
		dbi:  0,
	}

	dbdata.addInfo()
	dbdata.addDays()
	dbdata.addHours()
	dbdata.addTeachers()
	dbdata.addSubjects()
	dbdata.addRooms()
	dbdata.addRoomGroups()
	// RoomChoicesGroups: W365 has none of these – they must be generated
	// from the PreferredRooms lists of courses.
	dbdata.roomchoices = MapStrDb{}
	dbdata.data.RoomChoiceGroups = []db.RoomChoiceGroup{}
	dbdata.addGroups()
	dbdata.addClasses()
	dbdata.addCourses()
	dbdata.addCourses()
	dbdata.addSuperCourses()
	dbdata.addSubCourses()
	dbdata.addLessons()
	return dbdata.data
}

func (dbdata *xData) addSubjects() {
	dbdata.data.Subjects = []db.Subject{}
	dbdata.subjects = MapWDb{}
	dbdata.newsubjects = MapStrDb{}
	dbdata.subjectmap = MapWStr{}
	for _, d := range dbdata.w365.Subjects {
		sr := dbdata.nextId()
		dbdata.data.Subjects = append(dbdata.data.Subjects, db.Subject{
			Id:        sr,
			Tag:       d.Shortcut,
			Name:      d.Name,
			Reference: string(d.Id),
		})
		dbdata.subjects[d.Id] = sr
		dbdata.subjectmap[d.Id] = d.Shortcut
	}
}

func (dbdata *xData) addRooms() {
	dbdata.data.Rooms = []db.Room{}
	dbdata.rooms = MapWDb{}
	dbdata.roomtag = MapDbStr{}
	for _, d := range dbdata.w365.Rooms {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		rr := dbdata.nextId()
		dbdata.data.Rooms = append(dbdata.data.Rooms, db.Room{
			Id:           rr,
			Tag:          d.Shortcut,
			Name:         d.Name,
			NotAvailable: a,
			Reference:    string(d.Id),
		})
		dbdata.rooms[d.Id] = rr
		dbdata.roomtag[rr] = d.Shortcut
	}
}

func (dbdata *xData) addRoomGroups() {
	dbdata.data.RoomGroups = []db.RoomGroup{}
	dbdata.roomgroups = MapWDb{}
	for _, d := range dbdata.w365.RoomGroups {
		rlist := []DbRef{}
		for _, r := range d.Rooms {
			rr, ok := dbdata.rooms[r]
			if !ok {
				fmt.Printf("*ERROR* Unknown Room in RoomGroup %s:\n  %s\n",
					d.Id, r)
				continue
			}
			rlist = append(rlist, rr)
		}
		rr := dbdata.nextId()
		dbdata.data.RoomGroups = append(dbdata.data.RoomGroups, db.RoomGroup{
			Id:        rr,
			Tag:       d.Shortcut,
			Name:      d.Name,
			Reference: string(d.Id),
			Rooms:     rlist,
		})
		dbdata.roomgroups[d.Id] = rr
	}
}

func (dbdata *xData) addGroups() {
	// Every Group must be within one – and only one – Class Division.
	// To handle that, the data for the Groups is gathered here, but the
	// Elements are only added to the database when the Divisions are read.
	dbdata.data.Groups = []db.Group{}
	dbdata.pregroups = MapWStr{}
	dbdata.groups = MapWDb{}
	for _, d := range dbdata.w365.Groups {
		dbdata.pregroups[d.Id] = d.Shortcut
	}
}

func (dbdata *xData) addClasses() {
	dbdata.data.Classes = []db.Class{}
	dbdata.classes = MapWDb{}
	for _, d := range dbdata.w365.Classes {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		// Get the divisions and add their groups to the database.
		divs := []db.Division{}
		for _, wdiv := range d.Divisions {
			glist := []DbRef{}
			for _, g := range wdiv.Groups {
				gtag, ok := dbdata.pregroups[g] // get Tag
				if ok {
					// Add Group to database
					if _, nok := dbdata.groups[g]; nok {
						fmt.Printf("*ERROR* Group Defined in"+
							" multiple Divisions:\n  -- %s\n", g)
					}
					gr := dbdata.nextId()
					dbdata.data.Groups = append(dbdata.data.Groups, db.Group{
						Id:        gr,
						Tag:       gtag,
						Reference: string(g),
					})
					dbdata.groups[g] = gr
					glist = append(glist, gr)
				} else {
					fmt.Printf("*ERROR* Unknown Group in Class %s,"+
						" Division %s:\n  %s\n", d.Shortcut, wdiv.Name, g)
				}
			}
			divs = append(divs, db.Division{
				Name:      wdiv.Name,
				Groups:    glist,
				Reference: string(wdiv.Id),
			})
		}
		cr := dbdata.nextId()
		dbdata.data.Classes = append(dbdata.data.Classes, db.Class{
			Id:               cr,
			Name:             d.Name,
			Tag:              d.Shortcut,
			Level:            d.Level,
			Letter:           d.Letter,
			NotAvailable:     a,
			Divisions:        divs,
			MinLessonsPerDay: defaultMinus1(d.MinLessonsPerDay),
			MaxLessonsPerDay: defaultMinus1(d.MaxLessonsPerDay),
			MaxGapsPerDay:    defaultMinus1(d.MaxGapsPerDay),
			MaxGapsPerWeek:   defaultMinus1(d.MaxGapsPerWeek),
			MaxAfternoons:    defaultMinus1(d.MaxAfternoons),
			LunchBreak:       d.LunchBreak,
			ForceFirstHour:   d.ForceFirstHour,
			Reference:        string(d.Id),
		})
		dbdata.classes[d.Id] = cr
	}
}

func (dbdata *xData) readCourse(
	id W365Ref,
	subject W365Ref,
	subjects []W365Ref,
	groups []W365Ref,
	teachers []W365Ref,
	rooms []W365Ref,
) (DbRef, []DbRef, []DbRef, DbRef) {
	// Deal with subject
	var sr DbRef = 0
	var ok bool
	msg := "*ERROR* Course %s:\n  Unknown Subject: %s\n"
	if subject == "" {
		if len(subjects) == 1 {
			wsid := subjects[0]
			sr, ok = dbdata.subjects[wsid]
			if !ok {
				fmt.Printf(msg, id, wsid)
			}
		} else if len(subjects) > 1 {
			// Make a subject name
			sklist := []string{}
			for _, wsid := range subjects {
				// Need Shortcut field
				sk, ok := dbdata.subjectmap[wsid]
				if ok {
					sklist = append(sklist, sk)
				} else {
					fmt.Printf(msg, id, wsid)
				}
			}
			skname := strings.Join(sklist, ",")
			sr, ok = dbdata.newsubjects[skname]
			if !ok {
				sk := fmt.Sprintf("X%02d", len(dbdata.newsubjects)+1)
				sr = dbdata.nextId()
				dbdata.data.Subjects = append(dbdata.data.Subjects,
					db.Subject{
						Id:   sr,
						Tag:  sk,
						Name: skname,
					})
				dbdata.newsubjects[skname] = sr
			}
		}
	} else {
		if len(subjects) != 0 {
			fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n", id)
		}
		wsid := subject
		sr, ok = dbdata.subjects[wsid]
		if !ok {
			fmt.Printf(msg, id, wsid)
		}
	}
	// Deal with groups
	glist := []DbRef{}
	for _, g := range groups {
		gr, ok := dbdata.groups[g]
		// gr can refer to a Group or a Class.
		if !ok {
			// Check for class.
			gr, ok = dbdata.classes[g]
			if !ok {
				fmt.Printf("*ERROR* Unknown group in Course %s:\n  %s\n", id, g)
				continue
			}
		}
		glist = append(glist, gr)
	}
	// Deal with teachers
	tlist := []DbRef{}
	for _, t := range teachers {
		tr, ok := dbdata.teachers[t]
		if !ok {
			fmt.Printf("*ERROR* Unknown teacher in Course %s:\n  %s\n", id, t)
			continue
		}
		tlist = append(tlist, tr)
	}
	// Deal with rooms. W365 can have a single RoomGroup or a list of Rooms
	rclist := []DbRef{} // choice list
	var rm DbRef        // actual "room"
	for _, r := range rooms {
		rr, ok := dbdata.rooms[r]
		if ok {
			rclist = append(rclist, rr)
		} else {
			rm, ok = dbdata.roomgroups[r]
			if ok {
				if len(rooms) != 1 {
					rclist = []DbRef{}
					fmt.Printf(
						"*ERROR* Mixed Rooms and RoomGroups in Course %s\n", id)
				}
				break
			} else {
				fmt.Printf("*ERROR* Unknown room in Course %s:\n  %s\n", id, r)
				continue
			}
		}
	}
	if len(rclist) == 1 {
		// Take the single Room.
		rm = rclist[0]
	} else if len(rclist) > 1 {
		// Need a RoomChoiceGroup.
		// Reuse these if the same list appears again, but treat the
		// order as significant.
		rslist := []string{}
		for _, r := range rclist {
			rslist = append(rslist, dbdata.roomtag[r])
		}
		rs := strings.Join(rslist, ",")
		rm, ok = dbdata.roomchoices[rs]
		if !ok {
			rk := fmt.Sprintf("RC%03d", len(dbdata.roomchoices)+1)
			rm = dbdata.nextId()
			dbdata.data.RoomChoiceGroups = append(
				dbdata.data.RoomChoiceGroups, db.RoomChoiceGroup{
					Id:    rm,
					Tag:   rk,
					Name:  rs,
					Rooms: rclist,
				})
			dbdata.roomchoices[rs] = rm
		}
	}
	return sr, glist, tlist, rm
}

func (dbdata *xData) addCourses() {
	dbdata.data.Courses = []db.Course{}
	dbdata.courses = MapWDb{}
	for _, d := range dbdata.w365.Courses {
		sr, glist, tlist, rm := dbdata.readCourse(
			d.Id, d.Subject, d.Subjects, d.Groups, d.Teachers, d.PreferredRooms)
		cr := dbdata.nextId()
		dbdata.data.Courses = append(dbdata.data.Courses, db.Course{
			Id:        cr,
			Subject:   sr,
			Groups:    glist,
			Teachers:  tlist,
			Room:      rm,
			Reference: string(d.Id),
		})
		dbdata.courses[d.Id] = cr
	}
}

func (dbdata *xData) addSuperCourses() {
	dbdata.data.SuperCourses = []db.SuperCourse{}
	dbdata.supercourses = MapWDb{}
	for _, d := range dbdata.w365.SuperCourses {
		cr := dbdata.nextId()
		sr, ok := dbdata.subjects[d.Subject]
		if !ok {
			fmt.Printf("*ERROR* Unknown Subject in SuperCourse %s:\n  %s\n",
				d.Id, d.Subject)
			continue
		}
		dbdata.data.SuperCourses = append(dbdata.data.SuperCourses, db.SuperCourse{
			Id:        cr,
			Subject:   sr,
			Reference: string(d.Id),
		})
		dbdata.supercourses[d.Id] = cr
	}
}

func (dbdata *xData) addSubCourses() {
	dbdata.data.SubCourses = []db.SubCourse{}
	dbdata.subcourses = MapWDb{}
	for _, d := range dbdata.w365.SubCourses {
		sr, glist, tlist, rm := dbdata.readCourse(
			d.Id, d.Subject, d.Subjects, d.Groups, d.Teachers, d.PreferredRooms)
		sc, ok := dbdata.supercourses[d.SuperCourse]
		if !ok {
			fmt.Printf("*ERROR* Unknown SuperCourse in SubCourse %s:\n  %s\n",
				d.Id, d.SuperCourse)
			continue
		}
		cr := dbdata.nextId()
		dbdata.data.SubCourses = append(dbdata.data.SubCourses, db.SubCourse{
			Id:          cr,
			SuperCourse: sc,
			Subject:     sr,
			Groups:      glist,
			Teachers:    tlist,
			Room:        rm,
			Reference:   string(d.Id),
		})
		dbdata.subcourses[d.Id] = cr
	}
}

func (dbdata *xData) addLessons() {
	dbdata.data.Lessons = []db.Lesson{}
	for _, d := range dbdata.w365.Lessons {
		// The course can be either a Course or a SubCourse.
		crs, ok := dbdata.courses[d.Course]
		if !ok {
			crs, ok = dbdata.subcourses[d.Course]
			if !ok {
				fmt.Printf("*ERROR* Invalid course in Lesson %s:\n  -- %s\n",
					d.Id, d.Course)
				continue
			}
		}
		rlist := []DbRef{}
		for _, r := range d.LocalRooms {
			rr, ok := dbdata.rooms[r]
			if ok {
				rlist = append(rlist, rr)
			} else {
				fmt.Printf("*ERROR* Invalid room in Lesson %s:\n  -- %s\n",
					d.Id, r)
			}
		}
		dbdata.data.Lessons = append(dbdata.data.Lessons, db.Lesson{
			Id:        dbdata.nextId(),
			Course:    crs,
			Duration:  d.Duration,
			Day:       d.Day,
			Hour:      d.Hour,
			Fixed:     d.Fixed,
			Rooms:     rlist,
			Reference: string(d.Id),
		})
	}
}
*/