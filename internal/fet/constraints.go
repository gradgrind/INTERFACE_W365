package fet

import (
	"encoding/xml"
	"log"
)

type startingTime struct {
	XMLName            xml.Name `xml:"ConstraintActivityPreferredStartingTime"`
	Weight_Percentage  int
	Activity_Id        int
	Preferred_Day      string
	Preferred_Hour     string
	Permanently_Locked bool
	Active             bool
}

type minDaysBetweenActivities struct {
	XMLName                 xml.Name `xml:"ConstraintMinDaysBetweenActivities"`
	Weight_Percentage       int
	Consecutive_If_Same_Day bool
	Number_of_Activities    int
	Activity_Id             []int
	MinDays                 int
	Active                  bool
}

/* TODO
func gap_subject_activities(fetinfo *fetInfo,
	subject_activities []SubjectGroupActivities,
) {
	gsalist := []minDaysBetweenActivities{}
	for _, sga := range subject_activities {
		l := len(sga.Activities)
		// Adjust indexes for fet
		alist := []int{}
		// Skip if all activities are "fixed".
		allfixed := true
		for _, ai := range sga.Activities {
			alist = append(alist, ai+1)
			if !fetinfo.fixed_activities[ai] {
				allfixed = false
			}
		}
		if allfixed {
			continue
		}
		gsalist = append(gsalist, minDaysBetweenActivities{
			Weight_Percentage:       100,
			Consecutive_If_Same_Day: true,
			Number_of_Activities:    l,
			Activity_Id:             alist,
			MinDays:                 1,
			Active:                  true,
		})
	}
	// TODO
	//
	//	fetinfo.fetdata.Time_Constraints_List.ConstraintMinDaysBetweenActivities = gsalist
}
*/

type lunchBreak struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMaxHoursDailyInInterval"`
	Weight_Percentage   int
	Students            string
	Interval_Start_Hour string
	Interval_End_Hour   string
	Maximum_Hours_Daily int
	Active              bool
}

func lunch_break(
	fetinfo *fetInfo,
	lbconstraints *([]lunchBreak),
	cname string,
	lunchperiods []int,
) bool {
	// Assume the lunch periods are sorted, but not necessarily contiguous,
	// which is necessary for this constraint.
	lb1 := lunchperiods[0]
	lb2 := lunchperiods[len(lunchperiods)-1] + 1
	if lb2-lb1 != len(lunchperiods) {
		log.Printf(
			"\n=========================================\n"+
				"  !!!  INCOMPATIBLE DATA: lunch periods not contiguous,\n"+
				"       can't generate lunch-break constraint for class %s.\n"+
				"=========================================\n",
			cname)
		return false
	}
	lb := lunchBreak{
		Weight_Percentage:   100,
		Students:            cname,
		Interval_Start_Hour: fetinfo.hours[lb1],
		Interval_End_Hour:   fetinfo.hours[lb2],
		Maximum_Hours_Daily: len(lunchperiods) - 1,
		Active:              true,
	}
	*lbconstraints = append(*lbconstraints, lb)
	return true
}

type maxGapsPerWeek struct {
	XMLName           xml.Name `xml:"ConstraintStudentsSetMaxGapsPerWeek"`
	Weight_Percentage int
	Max_Gaps          int
	Students          string
	Active            bool
}

type minLessonsPerDay struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMinHoursDaily"`
	Weight_Percentage   int
	Minimum_Hours_Daily int
	Students            string
	Allow_Empty_Days    bool
	Active              bool
}

/*TODO:
For Teachers:
	MinLessonsPerDay interface{}
	MaxLessonsPerDay interface{}
	MaxDays          interface{}
	MaxGapsPerDay    interface{}
	MaxGapsPerWeek   interface{}
	MaxAfternoons    interface{}
	LunchBreak       bool
For Classes:
	MinLessonsPerDay interface{}
	MaxLessonsPerDay interface{}
	MaxGapsPerDay    interface{}
	MaxGapsPerWeek   interface{}
	MaxAfternoons    interface{}
	LunchBreak       bool
	ForceFirstHour   bool

Lunch breaks can be done using max-hours-in-interval constraint, but that
makes specification of max-gaps more difficult (becuase the lunch breaks
count as gaps).
The alternative is to add dummy lessons, clamped to the midday-break hours,
on the days where none of the midday-break hours are blocked. This can be a
problem if a class is finished earlier, but that may be a rare occurrence.

The different-days constraint for lessons belonging to a single course can
be added automatically, but it should be posible to disable it by passing in
an appropriate constraint. Thus, the built-in constraint must be traceable.
There could be a separate constraint to link different courses – the
alternative being a subject/atomic-group search.
*/

func addDifferentDaysConstraints(fetinfo *fetInfo) {
	mdba := []minDaysBetweenActivities{}
	for cref, cinfo := range fetinfo.courseInfo {
		if len(cinfo.activities) < 2 {
			continue
		}
		// Need the Acivity_Ids for the Lessons, and whether they are fixed.
		// No two fixed activities should be different-dayed.

		fixeds := []int{}
		unfixeds := []int{}
		for i, l := range cinfo.lessons {
			if l.Fixed {
				fixeds = append(fixeds, cinfo.activities[i])
			} else {
				unfixeds = append(unfixeds, cinfo.activities[i])
			}
		}

		if len(fixeds) <= 1 {
			fetinfo.differentDayConstraints[cref] = []int{len(mdba)}
			mdba = append(mdba, minDaysBetweenActivities{
				Weight_Percentage:       100,
				Consecutive_If_Same_Day: true,
				Number_of_Activities:    len(unfixeds),
				Activity_Id:             cinfo.activities,
				MinDays:                 1,
				Active:                  true,
			})
			continue
		}

		if len(unfixeds) == 0 {
			continue
		}

		ddc := []int{} // Collect indexes within mdba
		for _, aid := range fixeds {
			aids := []int{aid}
			aids = append(aids, unfixeds...)
			ddc = append(ddc, len(mdba))
			mdba = append(mdba, minDaysBetweenActivities{
				Weight_Percentage:       100,
				Consecutive_If_Same_Day: true,
				Number_of_Activities:    len(aids),
				Activity_Id:             aids,
				MinDays:                 1,
				Active:                  true,
			})
		}
		fetinfo.differentDayConstraints[cref] = ddc
	}
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintMinDaysBetweenActivities = mdba
}
