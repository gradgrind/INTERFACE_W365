package fet

import (
	"encoding/xml"
	"fmt"
)

// const GROUP_SEP = ","
// const DIV_SEP = "|"
const CLASS_GROUP_SEP = "."
const ATOMIC_GROUP_SEP1 = "#"
const ATOMIC_GROUP_SEP2 = "~"

type fetCategory struct {
	//XMLName             xml.Name `xml:"Category"`
	Number_of_Divisions int
	Division            []string
}

type fetSubgroup struct {
	Name string // 13.m.MaE
	//Number_of_Students int // 0
	//Comments string // ""
}

type fetGroup struct {
	Name string // 13.K
	//Number_of_Students int // 0
	//Comments string // ""
	Subgroup []fetSubgroup
}

type fetClass struct {
	//XMLName  xml.Name `xml:"Year"`
	Name      string
	Long_Name string
	Comments  string
	//Number_of_Students int (=0)
	// The information regarding categories, divisions of each category,
	// and separator is only used in the dialog to divide the year
	// automatically by categories.
	Number_of_Categories int
	Separator            string // CLASS_GROUP_SEP
	Category             []fetCategory
	Group                []fetGroup
}

type fetStudentsList struct {
	XMLName xml.Name `xml:"Students_List"`
	Year    []fetClass
}

type studentsNotAvailable struct {
	XMLName                       xml.Name `xml:"ConstraintStudentsSetNotAvailableTimes"`
	Weight_Percentage             int
	Students                      string
	Number_of_Not_Available_Times int
	Not_Available_Time            []notAvailableTime
	Active                        bool
}

//TODO

func getClasses(fetinfo *fetInfo) {
	//items := []fetClass{}
	//natimes := []studentsNotAvailable{}
	//lunchperiods := fetinfo.db.Info.MiddayBreak
	//lunchconstraints := []lunchBreak{}
	//maxgaps := []maxGapsPerWeek{}
	//minlessons := []minLessonsPerDay{}
	for _, cl := range fetinfo.db.Classes {
		//TODO
		//cgs := fetinfo.wzdb.AtomicGroups.Class_Groups[c]
		//agmap := fetinfo.wzdb.AtomicGroups.Group_Atomics
		//cags := agmap[wzbase.ClassGroup{
		//	CIX: c, GIX: 0,
		//}]

		//		divs := cl.DIVISIONS
		//nc := 0
		//		if len(divs) > 0 {
		//if cags.GetCardinality() > 1 {
		//	nc = 1
		//}
		//calt := cl.SORTING //?
		cname := cl.Tag
		clAGs := fetinfo.atomicGroups[cl.Id]
		fmt.Printf("##### cags %s: %+v\n", cname, clAGs)

	}
	return
	/*
		for {

			groups := []fetGroup{}
			if cags.GetCardinality() > 1 {
				for _, cg := range cgs {
					g := fetinfo.ref2fet[cg.GIX]
					gags := agmap[cg]
					subgroups := []fetSubgroup{}
					for _, ag := range gags.ToArray() {
						subgroups = append(subgroups,
							fetSubgroup{Name: fmt.Sprintf("%s.%03d", cname, ag)},
						)
						//ag_gs[int(ag)] = append(ag_gs[int(ag)], g)
					}
					groups = append(groups, fetGroup{
						Name:     fmt.Sprintf("%s.%s", cname, g),
						Subgroup: subgroups,
					})
				}
			}

			// Use the Comments field as an additional specification of the
			// partitioning.
			slcum := []string{}
			active_divisions := fetinfo.wzdb.ActiveDivisions[c]
			categories := []fetCategory{}
			for _, divl := range active_divisions {
				strcum := []string{}
				for _, i := range divl {
					strcum = append(strcum, fetinfo.ref2fet[i])
				}
				categories = append(categories, fetCategory{
					Number_of_Divisions: len(divl),
					Division:            strcum,
				})
				slcum = append(slcum, strings.Join(strcum, GROUP_SEP))
			}
			strdivs := strings.Join(slcum, DIV_SEP)
			//fmt.Printf("??? ActiveDivisions %s (%s): %+v\n",
			//	cname, cl.SORTING, strdivs)
			items = append(items, fetClass{
				Name:                 cname,
				Long_Name:            cl.NAME,
				Comments:             strdivs,
				Separator:            ".",
				Number_of_Categories: len(categories),
				Category:             categories,
				Group:                groups,
			})

			//fmt.Printf("\nCLASS %s: %+v\n", cl.SORTING, cl.DIVISIONS)

			// ************************************************************
			// The following constraints don't concern dummy classes ending
			// in "X".
			if strings.HasSuffix(cname, "X") {
				continue
			}

			// "Not available" times
			// Seek also the days where a lunch-break is necessary – those days
			// where none of the lunch-break periods are blocked.
			lbdays := []int{}
			nats := []notAvailableTime{}
			for d, dna := range cl.NOT_AVAILABLE {
				lbd := true
				for _, h := range dna {
					if lbd {
						for _, hh := range lunchperiods {
							if hh == h {
								lbd = false
								break
							}
						}
					}
					nats = append(nats,
						notAvailableTime{
							Day: fetinfo.days[d], Hour: fetinfo.hours[h]})
				}
				if lbd {
					lbdays = append(lbdays, d)
				}
			}

			if len(nats) > 0 {
				natimes = append(natimes,
					studentsNotAvailable{
						Weight_Percentage:             100,
						Students:                      cname,
						Number_of_Not_Available_Times: len(nats),
						Not_Available_Time:            nats,
						Active:                        true,
					})
			}
			//fmt.Printf("==== %s: %+v\n", cname, nats)

			// Limit gaps on a weekly basis.
			mgpw := 0 //TODO: An additional tweak may be needed for some classes.
			// Handle lunch breaks: The current approach counts lunch breaks as
			// gaps, so the gaps-per-week must be adjusted accordingly.
			if len(lbdays) > 0 {
				// Need lunch break(s).
				// This uses a general "max-lessons-in-interval" constraint.
				// As an alternative, adding dummy lessons (with time constraint)
				// can offer some advantages, like easing gap handling.
				// Set max-gaps-per-week accordingly.
				if lunch_break(fetinfo, &lunchconstraints, cname, lunchperiods) {
					mgpw += len(lbdays)
				}
			}
			// Add the gaps constraint.
			maxgaps = append(maxgaps, maxGapsPerWeek{
				Weight_Percentage: 100,
				Max_Gaps:          mgpw,
				Students:          cname,
				Active:            true,
			})

			// Minimum lessons per day
			mlpd0 := cl.CONSTRAINTS["MinLessonsPerDay"]
			mlpd, err := strconv.Atoi(mlpd0)
			if err != nil {
				log.Fatalf("INVALID MinLessonsPerDay: %s // %v\n", mlpd0, err)
			}
			minlessons = append(minlessons, minLessonsPerDay{
				Weight_Percentage:   100,
				Minimum_Hours_Daily: mlpd,
				Students:            cname,
				Allow_Empty_Days:    false,
				Active:              true,
			})
		}
		fetinfo.fetdata.Students_List = fetStudentsList{
			Year: items,
		}
		fetinfo.fetdata.Time_Constraints_List.
			ConstraintStudentsSetNotAvailableTimes = natimes
		fetinfo.fetdata.Time_Constraints_List.
			ConstraintStudentsSetMaxHoursDailyInInterval = lunchconstraints
		fetinfo.fetdata.Time_Constraints_List.
			ConstraintStudentsSetMaxGapsPerWeek = maxgaps
		fetinfo.fetdata.Time_Constraints_List.
			ConstraintStudentsSetMinHoursDaily = minlessons
	*/
}
