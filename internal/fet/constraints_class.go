package fet

import "log"

/*TODO:
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

/*TODO:
For Teachers:
	LunchBreak       bool

Lunch-breaks can be done using max-hours-in-interval constraint, but that
makes specification of max-gaps more difficult (becuase the lunch breaks
count as gaps).
The alternative is to add dummy lessons, clamped to the midday-break hours,
on the days where none of the midday-break hours are blocked. This can be a
problem if a class is finished earlier, but that may be a rare occurrence.


func teacherLunchBreaks(fetinfo *fetInfo, t *w365tt.Teacher) {
	mbhours := fetinfo.db.Info.MiddayBreak
	days := make([]bool, len(fetinfo.days))
	d := 0
	for _, ts := range t.NotAvailable {
		if ts.Day < d {
			continue
		}
		if slices.Contains(mbhours, ts.Hour) {
			days[ts.Day] = true
			d = ts.Day + 1
		}
	}

	// Possibility 1: Add a dummy lesson for each day in days, constrained
	// to the hours in mbhours.

	// Possibility 2: Add a lunch-break constraint. This doesn't need the
	// days, but they may be useful for adjusting the max gaps?

	// Add dummy lessons for lunch-breaks.
	for d, nok := range days {
		if !nok {
			//TODO
			//addTeacherLunchBreak(fetinfo, t.Tag, d)
			aid := addTeacherLunchBreak(fetinfo, t.Tag, d)
			fmt.Printf("§LB (%d): %s / %d\n", aid, t.Tag, d)
		}
	}
}

func addTeacherLunchBreak(fetinfo *fetInfo, ttag string, day int) int {
	acl := &fetinfo.fetdata.Activities_List
	aid := len(acl.Activity) + 1
	acl.Activity = append(acl.Activity, fetActivity{
		Id:                aid,
		Teacher:           []string{ttag},
		Subject:           LUNCH_BREAK_TAG,
		Students:          []string{},
		Active:            true,
		Total_Duration:    1,
		Duration:          1,
		Activity_Group_Id: 0,
	})
	return aid
}

/*

<ConstraintActivitiesPreferredTimeSlots>
	<Weight_Percentage>100</Weight_Percentage>
	<Teacher></Teacher>
	<Students></Students>
	<Subject>-lb-</Subject>
	<Activity_Tag></Activity_Tag>
	<Duration></Duration>
	<Number_of_Preferred_Time_Slots>2</Number_of_Preferred_Time_Slots>
	<Preferred_Time_Slot>
		<Preferred_Day>Mo</Preferred_Day>
		<Preferred_Hour>(6)</Preferred_Hour>
	</Preferred_Time_Slot>
	<Preferred_Time_Slot>
		<Preferred_Day>Mo</Preferred_Day>
		<Preferred_Hour>(7)</Preferred_Hour>
	</Preferred_Time_Slot>
	<Active>true</Active>
	<Comments></Comments>
</ConstraintActivitiesPreferredTimeSlots>

*/
