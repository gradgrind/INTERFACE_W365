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
