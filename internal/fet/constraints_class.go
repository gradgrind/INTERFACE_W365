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
