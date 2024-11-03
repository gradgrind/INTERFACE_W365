package fet

/*TODO:
For Teachers:
	MinLessonsPerDay interface{}
	MaxLessonsPerDay interface{}
	MaxDays          interface{}
	MaxGapsPerDay    interface{}
	MaxGapsPerWeek   interface{}
	MaxAfternoons    interface{}
	LunchBreak       bool

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

// TODO: Add other constraints

func addTeacherConstraints(fetinfo *fetInfo) {
	tmaxdpw := []maxDaysT{}
	tminlpd := []minLessonsPerDayT{}
	tmaxlpd := []maxLessonsPerDayT{}
	for _, t := range fetinfo.db.Teachers {
		tmaxdpw = append(tmaxdpw, maxDaysT{
			Weight_Percentage: 100,
			Teacher:           t.Tag,
			Max_Days_Per_Week: t.MaxDays.(int),
			Active:            true,
		})

		tminlpd = append(tminlpd, minLessonsPerDayT{
			Weight_Percentage:   100,
			Teacher:             t.Tag,
			Minimum_Hours_Daily: t.MinLessonsPerDay.(int),
			Allow_Empty_Days:    false,
			Active:              true,
		})

		tmaxlpd = append(tmaxlpd, maxLessonsPerDayT{
			Weight_Percentage:   100,
			Teacher:             t.Tag,
			Maximum_Hours_Daily: t.MaxLessonsPerDay.(int),
			Active:              true,
		})

	}
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMaxDaysPerWeek = tmaxdpw
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMinHoursDaily = tminlpd
}
