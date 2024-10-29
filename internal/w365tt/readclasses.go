package w365tt

import (
	"fmt"
	"log"
)

func (dbdata *xData) readClasses() {
	// Every Class-Group must be within one – and only one – Class-Division.
	// To handle that, the data for the Groups is first gathered here, but
	// the Elements are only added to the database map when the Divisions
	// are read.
	pregroups := map[Ref]*Group{}
	for i, n := range dbdata.data.Groups {
		pregroups[n.Id] = &dbdata.data.Groups[i]
	}

	for i := 0; i < len(dbdata.data.Classes); i++ {
		n := &dbdata.data.Classes[i]
		dbdata.elements[n.Id] = n

		if len(n.NotAvailable) == 0 {
			// Avoid a null value
			n.NotAvailable = []TimeSlot{}
		}
		if n.MinLessonsPerDay == nil {
			n.MinLessonsPerDay = -1
		}
		if n.MaxLessonsPerDay == nil {
			n.MaxLessonsPerDay = -1
		}
		if n.MaxGapsPerDay == nil {
			n.MaxGapsPerDay = -1
		}
		if n.MaxGapsPerWeek == nil {
			n.MaxGapsPerWeek = -1
		}
		if n.MaxAfternoons == nil {
			n.MaxAfternoons = -1
		}

		// Get the divisions and add their groups to the database.
		for i, wdiv := range n.Divisions {
			glist := []Ref{}
			for _, g := range wdiv.Groups {
				// get Tag
				group, ok := pregroups[g]
				if ok {
					// Add Group to database, if it's not already there
					if _, nok := dbdata.elements[g]; nok {
						log.Fatalf("*ERROR* Group Defined in"+
							" multiple Divisions:\n  -- %s\n", g)
					}
					dbdata.elements[g] = group
					glist = append(glist, g)
				} else {
					fmt.Printf("*ERROR* Unknown Group in Class %s,"+
						" Division %s:\n  %s\n", n.Tag, wdiv.Name, g)
				}
			}
			// Accept Divisions which have too few Groups at this stage.
			if len(glist) < 2 {
				fmt.Printf("*WARNING* In Class %s,"+
					" not enough valid Groups (>1) in Division %s\n",
					n.Tag, wdiv.Name)
			}
			n.Divisions[i].Groups = glist
		}
	}
}
