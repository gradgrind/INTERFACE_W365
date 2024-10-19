package w365_tt

import (
	"fmt"
	"strings"
)

func test_ids_exist(w365 *W365TT, idmap IdMap) {
	for _, c := range w365.Teachers {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Teachers.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Teachers.Categories %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Rooms {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Rooms.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Rooms.Categories %s\n", x)
				}
			}
		}
		if len(c.RoomGroups) != 0 {
			for _, x := range strings.Split(c.RoomGroups, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Rooms.RoomGroups %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Classes {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Classes.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Classes.Categories %s\n", x)
				}
			}
		}
		if len(c.Divisions) != 0 {
			for _, x := range strings.Split(c.Divisions, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Classes.Divisions %s\n", x)
				}
			}
		}
		if len(c.Groups) != 0 {
			for _, x := range strings.Split(c.Groups, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Classes.Groups %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Groups {
		gstr, ok := groupTagFull(idmap, c.Id)
		if !ok {
			fmt.Printf("  !!! Invalid Group: %s\n", gstr)
		}
	}
	for _, c := range w365.Divisions {
		if len(c.Groups) != 0 {
			for _, x := range strings.Split(c.Groups, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Divisions.Groups %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Courses {
		if len(c.Subjects) != 0 {
			for _, x := range strings.Split(c.Subjects, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Courses.Subjects %s\n", x)
				}
			}
		}
		if len(c.Groups) != 0 {
			for _, x := range strings.Split(c.Groups, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Courses.Groups %s\n", x)
				}
			}
		}
		if len(c.Teachers) != 0 {
			for _, x := range strings.Split(c.Teachers, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Courses.Teachers %s\n", x)
				}
			}
		}
		if len(c.PreferredRooms) != 0 {
			for _, x := range strings.Split(c.PreferredRooms, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Courses.PreferredRooms %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Courses.Categories %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.EpochPlanCourses {
		if len(c.Subjects) != 0 {
			for _, x := range strings.Split(c.Subjects, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! EpochPlanCourses.Subjects %s\n", x)
				}
			}
		}
		if len(c.Groups) != 0 {
			for _, x := range strings.Split(c.Groups, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! EpochPlanCourses.Groups %s\n", x)
				}
			}
		}
		if len(c.Teachers) != 0 {
			for _, x := range strings.Split(c.Teachers, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! EpochPlanCourses.Teachers %s\n", x)
				}
			}
		}
		if len(c.PreferredRooms) != 0 {
			for _, x := range strings.Split(c.PreferredRooms, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! EpochPlanCourses.PreferredRooms %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! EpochPlanCourses.Categories %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Lessons {
		if len(c.Course) == 0 {
			// Let's assume this Lesson is invalid/unused
			continue
			//fmt.Println("  !!! Lesson with no Course")
		} else {
			crs, ok := idmap.Id2Node[c.Course]
			if !ok {
				// Let's assume this Lesson is invalid/unused
				continue
				//fmt.Printf("  !!! Lessons.Course %s\n", c.Course)
			} else {
				_, ok := crs.(*Course)
				if ok {
					//fmt.Printf("  +++ Lesson : Course %s : %s\n", c.Id, c.Course)
				} else {
					_, ok := crs.(*EpochPlanCourse)
					if ok {
						fmt.Printf("  +++ Lesson : EpochCourse %s : %s\n", c.Id, c.Course)
					} else {
						fmt.Printf("  !!! Lessons.Course? %s\n", c.Course)
					}
				}
			}
		}
		if len(c.LocalRooms) != 0 {
			for _, x := range strings.Split(c.LocalRooms, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Lessons.LocalRooms %s\n", x)
				}
			}
		}
		if len(c.EpochPlan) != 0 {
			for _, x := range strings.Split(c.EpochPlan, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Lessons.EpochPlan %s\n", x)
				}
			}
		}
		if len(c.EpochPlanGrade) != 0 {
			for _, x := range strings.Split(c.EpochPlanGrade, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Lessons.EpochPlanGrade %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Fractions {
		if len(c.SuperGroups) != 0 {
			for _, x := range strings.Split(c.SuperGroups, ",") {
				_, ok := idmap.Id2Node[x]
				if !ok {
					fmt.Printf("  !!! Fractions.SuperGroups %s\n", x)
				}
			}
		}
	}

}
