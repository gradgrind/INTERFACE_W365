package fet

import (
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"slices"
)

/*TODO:
For Teachers:
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

// TODO
func teacherLunchBreaks(fetinfo *fetInfo, t *w365tt.Teacher) {
	mbhours := fetinfo.db.Info.MiddayBreak
	days := []int{}
	d := 0
	for _, ts := range t.NotAvailable {
		if ts.Day < d {
			continue
		}
		if slices.Contains(mbhours, ts.Hour) {
			days = append(days, ts.Day)
			d = ts.Day + 1
		}
	}

	// Possibility 1: Add a dummy lesson for each day in days, constrained
	// to the hours in mbhours.

	// Possibility 2: Add a lunch-break constraint. This doesn't need the
	// days, but they may be useful for adjusting the max gaps?
}

func addTeacherConstraints(fetinfo *fetInfo) {
	tmaxdpw := []maxDaysT{}
	tminlpd := []minLessonsPerDayT{}
	tmaxlpd := []maxLessonsPerDayT{}
	tmaxgpd := []maxGapsPerDayT{}
	tmaxgpw := []maxGapsPerWeekT{}
	tmaxaft := []maxDaysinIntervalPerWeekT{}
	ndays := len(fetinfo.days)
	nhours := len(fetinfo.hours)

	for _, t := range fetinfo.db.Teachers {
		n := int(t.MaxDays.(float64))
		if n >= 0 && n < ndays {
			tmaxdpw = append(tmaxdpw, maxDaysT{
				Weight_Percentage: 100,
				Teacher:           t.Tag,
				Max_Days_Per_Week: n,
				Active:            true,
			})
		}

		n = int(t.MinLessonsPerDay.(float64))
		if n >= 0 && n <= nhours {
			tminlpd = append(tminlpd, minLessonsPerDayT{
				Weight_Percentage:   100,
				Teacher:             t.Tag,
				Minimum_Hours_Daily: n,
				Allow_Empty_Days:    false,
				Active:              true,
			})
		}

		n = int(t.MaxLessonsPerDay.(float64))
		if n >= 0 && n < nhours {
			tmaxlpd = append(tmaxlpd, maxLessonsPerDayT{
				Weight_Percentage:   100,
				Teacher:             t.Tag,
				Maximum_Hours_Daily: n,
				Active:              true,
			})
		}

		n = int(t.MaxGapsPerDay.(float64))
		if n >= 0 {
			tmaxgpd = append(tmaxgpd, maxGapsPerDayT{
				Weight_Percentage: 100,
				Teacher:           t.Tag,
				Max_Gaps:          n,
				Active:            true,
			})
		}

		n = int(t.MaxGapsPerWeek.(float64))
		if n >= 0 {
			tmaxgpw = append(tmaxgpw, maxGapsPerWeekT{
				Weight_Percentage: 100,
				Teacher:           t.Tag,
				Max_Gaps:          n,
				Active:            true,
			})
		}

		i := fetinfo.db.Info.FirstAfternoonHour
		n = int(t.MaxAfternoons.(float64))
		if n >= 0 && i > 0 {
			tmaxaft = append(tmaxaft, maxDaysinIntervalPerWeekT{
				Weight_Percentage:   100,
				Teacher:             t.Tag,
				Interval_Start_Hour: fetinfo.hours[i],
				Interval_End_Hour:   "", // end of day
				Max_Days_Per_Week:   n,
				Active:              true,
			})
		}

	}
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMaxDaysPerWeek = tmaxdpw
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMinHoursDaily = tminlpd
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMaxHoursDaily = tmaxlpd
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMaxGapsPerDay = tmaxgpd
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherMaxGapsPerWeek = tmaxgpw
	fetinfo.fetdata.Time_Constraints_List.
		ConstraintTeacherIntervalMaxDaysPerWeek = tmaxaft
}
