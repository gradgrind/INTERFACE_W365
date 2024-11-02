package fet

import (
	"encoding/xml"
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"log"
	"slices"
)

type fetActivity struct {
	XMLName           xml.Name `xml:"Activity"`
	Id                int
	Teacher           []string
	Subject           string
	Activity_Tag      string `xml:",omitempty"`
	Students          []string
	Active            bool
	Total_Duration    int
	Duration          int
	Activity_Group_Id int
	Comments          string
}

type fetActivitiesList struct {
	XMLName  xml.Name `xml:"Activities_List"`
	Activity []fetActivity
}

type fetActivityTag struct {
	XMLName   xml.Name `xml:"Activity_Tag"`
	Name      string
	Printable bool
}

type fetActivityTags struct {
	XMLName      xml.Name `xml:"Activity_Tags_List"`
	Activity_Tag []fetActivityTag
}

func gatherCourseInfo(fetinfo *fetInfo) {
	// Gather the Groups, Teachers and "rooms" for the Courses and
	// SuperCourses with lessons (only).
	// Gather the Lessons for these Courses and SuperCourses.
	// Also, the SuperCourses (with lessons) get a list of their
	// SubCourses.
	db := fetinfo.db
	fetinfo.superSubs = make(map[Ref][]Ref)
	fetinfo.courseInfo = make(map[Ref]courseInfo)
	roomData := map[Ref][]Ref{}
	// Collect Courses with Lessons.
	for _, l := range db.Lessons {
		lcref := l.Course
		cinfo, ok := fetinfo.courseInfo[lcref]
		if ok {
			// If the course has already been handled, just add the lesson.
			cinfo.lessons = append(cinfo.lessons, l.Id)
			fetinfo.courseInfo[lcref] = cinfo
			continue
		}
		// First encounter with the course.
		var subject Ref
		var groups []Ref
		var teachers []Ref
		var rooms []Ref
		lessons := []Ref{l.Id}

		c := db.Elements[lcref] // can be Course or SuperCourse
		cnode, ok := c.(*w365tt.Course)
		if ok {
			subject = cnode.Subject
			groups = cnode.Groups
			teachers = cnode.Teachers
			rooms = []Ref{cnode.Room}
		} else {
			spc, ok := c.(*w365tt.SuperCourse)
			if !ok {
				log.Fatalf(
					"*ERROR* Invalid Course in Lesson %s:\n  %s\n",
					l.Id, lcref)
			}
			fetinfo.superSubs[lcref] = []Ref{}
			subject = spc.Subject
			groups = []Ref{}
			teachers = []Ref{}
			rooms = []Ref{}
		}
		fetinfo.courseInfo[lcref] = courseInfo{
			subject:  subject,
			groups:   groups,
			teachers: teachers,
			//rooms: filled later
			lessons: lessons,
		}
		roomData[lcref] = rooms
	}
	// Now find the SubCourses
	for _, sbc := range db.SubCourses {
		spc := sbc.SuperCourse
		cinfo, ok := fetinfo.courseInfo[spc]
		if ok {
			// Only fill SuperCourses which have Lessons
			fetinfo.superSubs[spc] = append(fetinfo.superSubs[spc], sbc.Id)
			// Add groups
			cglist := append(cinfo.groups, sbc.Groups...)
			slices.Sort(cglist)
			cglist = slices.Compact(cglist)
			cinfo.groups = make([]Ref, len(cglist))
			copy(cinfo.groups, cglist)
			// Add teachers
			ctlist := append(cinfo.teachers, sbc.Teachers...)
			slices.Sort(ctlist)
			ctlist = slices.Compact(ctlist)
			cinfo.teachers = make([]Ref, len(ctlist))
			copy(cinfo.teachers, ctlist)
			// Add rooms
			crlist := append(roomData[spc], sbc.Room)
			slices.Sort(crlist)
			crlist = slices.Compact(crlist)
			roomData[spc] = crlist

			fetinfo.courseInfo[spc] = cinfo
		}
	}
	// Prepare the internal room structure, filtering the room lists of
	// the SuperCourses.
	for cref, crlist := range roomData {
		// TODO
		// Join all Rooms and the Rooms from RoomGroups into a "compulsory"
		// list. Then go through the RoomChoiceGroups. If one contains a
		// compulsory room, ignore the choice.
		// The result is a list of Rooms and a list of room-choice-lists,
		// which can be converted into a fet virtual room.
		rooms := []Ref{}
		roomChoices := [][]Ref{}
		for _, rref := range crlist {
			rx := fetinfo.db.Elements[rref]
			_, ok := rx.(*w365tt.Room)
			if ok {
				rooms = append(rooms, rref)
			} else {
				rg, ok := rx.(*w365tt.RoomGroup)
				if ok {
					rooms = append(rooms, rg.Rooms...)
				} else {
					rc, ok := rx.(*w365tt.RoomChoiceGroup)
					if !ok {
						log.Fatalf(
							"*BUG* Invalid room in course %s:\n  %s\n",
							cref, rref)
					}
					roomChoices = append(roomChoices, rc.Rooms)
				}
			}
		}
		// Remove duplicates in Room list.
		slices.Sort(rooms)
		rooms = slices.Compact(rooms)
		// Filter choice lists.
		roomChoices = slices.DeleteFunc(roomChoices, func(rcl []Ref) bool {
			for _, rc := range rcl {
				if slices.Contains(rooms, rc) {
					return true
				}
			}
			return false
		})
		cinfo := fetinfo.courseInfo[cref]
		cinfo.room = virtualRoom{
			rooms:       rooms,
			roomChoices: roomChoices,
		}
		fetinfo.courseInfo[cref] = cinfo
	}
}

// Generate the fet activties.
func getActivities(fetinfo *fetInfo) {

	// ************* Start with the activity tags
	tags := []fetActivityTag{}
	/* ???
	s2tag := map[string]string{}
	for _, ts := range tagged_subjects {
		tag := fmt.Sprintf("Tag_%s", ts)
		s2tag[ts] = tag
		tags = append(tags, fetActivityTag{
			Name: tag,
		})
	}
	*/
	fetinfo.fetdata.Activity_Tags_List = fetActivityTags{
		Activity_Tag: tags,
	}
	// ************* Now the activities
	activities := []fetActivity{}
	lessonList := []*w365tt.Lesson{{}} // with empty first entry, because
	// Activity Ids start at 1
	aid := 0
	for _, cinfo := range fetinfo.courseInfo {
		// Teachers
		tlist := []string{}
		for _, ti := range cinfo.teachers {
			tlist = append(tlist, fetinfo.ref2fet[ti])
		}
		// Groups
		glist := []string{}
		for _, cgref := range cinfo.groups {
			glist = append(glist, fetinfo.ref2fet[cgref])
		}
		/* ???
		atag := ""
		if slices.Contains(tagged_subjects, sbj) {
			atag = fmt.Sprintf("Tag_%s", sbj)
		}
		*/

		// Generate the Activities for this course (one per Lesson).
		agid := 0 // first activity should have Id = 1
		if len(cinfo.lessons) > 1 {
			agid = aid + 1
		}
		totalDuration := 0
		llist := []*w365tt.Lesson{}
		for _, lref := range cinfo.lessons {
			l := fetinfo.db.Elements[lref].(*w365tt.Lesson)
			llist = append(llist, l)
			totalDuration += l.Duration
		}
		for _, l := range llist {
			aid++
			activities = append(activities,
				fetActivity{
					Id:       aid,
					Teacher:  tlist,
					Subject:  fetinfo.ref2fet[cinfo.subject],
					Students: glist,
					//Activity_Tag:      atag,
					Active:            true,
					Total_Duration:    totalDuration,
					Duration:          l.Duration,
					Activity_Group_Id: agid,
					Comments:          string(l.Id),
				},
			)
			lessonList = append(lessonList, l)
		}
	}
	fetinfo.fetdata.Activities_List = fetActivitiesList{
		Activity: activities,
	}

	/*TODO: Deal with room constraints

	fixed_rooms := []fixedRoom{}
	room_choices := []roomChoice{}
	virtual_rooms := map[string]string{}

			addRoomConstraint(fetinfo,
			&fixed_rooms,
			&room_choices,
			virtual_rooms,
			acts,
			activity.RoomNeeds,
		)
	}
	// Now generate the full list of fet activities
	starting_times := []startingTime{}
	items := []fetActivity{}
	fetinfo.fixed_activities = make([]bool, len(activities))
	for i, activity := range activities {
		ci := activity.Course
		fetact := course_act[ci]
		fetact.Id = i + 1
		fetact.Duration = activity.Duration
		items = append(items, fetact)
		// Activity placement
		day := activity.Day
		if day >= 0 {
			hour := activity.Hour
			starting_times = append(starting_times, startingTime{
				Weight_Percentage:  100,
				Activity_Id:        i + 1,
				Preferred_Day:      fetinfo.days[day],
				Preferred_Hour:     fetinfo.hours[hour],
				Permanently_Locked: true,
				Active:             true,
			})
			fetinfo.fixed_activities[i] = true
		}
	}
	fetinfo.fetdata.Activities_List = fetActivitiesList{
		Activity: items,
	}
	fetinfo.fetdata.Time_Constraints_List.ConstraintActivityPreferredStartingTime = starting_times
	fetinfo.fetdata.Space_Constraints_List.ConstraintActivityPreferredRoom = fixed_rooms
	fetinfo.fetdata.Space_Constraints_List.ConstraintActivityPreferredRooms = room_choices
	*/
}
