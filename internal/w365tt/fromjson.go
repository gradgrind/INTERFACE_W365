package w365tt

import (
	"encoding/json"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/db"
	"io"
	"log"
	"os"
	"strings"
)

func ReadJSON(jsonpath string) W365TopLevel {
	// Open the  JSON file
	jsonFile, err := os.Open(jsonpath)
	if err != nil {
		log.Fatal(err)
	}
	// Remember to close the file at the end of the function
	defer jsonFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	log.Printf("*+ Reading: %s\n", jsonpath)
	v := W365TopLevel{}
	err = json.Unmarshal(byteValue, &v)
	if err != nil {
		log.Fatalf("Could not unmarshal json: %s\n", err)
	}

	DeMultipleSubjects(&v)

	return v
}

type xData struct {
	w365  W365TopLevel
	data  db.DbTopLevel
	dbi   db.DbRef
	idmap map[W365Ref]db.DbRef
}

func LoadJSON(jsonpath string) db.DbTopLevel {
	dbdata := xData{
		w365:  ReadJSON(jsonpath),
		data:  db.DbTopLevel{},
		dbi:   0,
		idmap: map[W365Ref]db.DbRef{},
	}

	dbdata.addInfo()
	dbdata.addDays()
	dbdata.addHours()

	return dbdata.data
}

func (dbdata *xData) nextId(w365Id W365Ref) db.DbRef {
	dbdata.dbi++
	dbdata.idmap[w365Id] = dbdata.dbi
	return dbdata.dbi
}

func (dbdata *xData) addInfo() {
	dbdata.data.Info = db.Info{
		FirstAfternoonHour: dbdata.w365.W365TT.FirstAfternoonHour,
		MiddayBreak:        dbdata.w365.W365TT.MiddayBreak,
	}
}

func (dbdata *xData) addDays() {
	for _, d := range dbdata.w365.Days {
		dbdata.data.Days = append(dbdata.data.Days, db.Day{
			Id:   dbdata.nextId(d.Id),
			Type: db.TypeDAY,
			Tag:  d.Shortcut,
			Name: d.Name,
		})
	}
}

func (dbdata *xData) addHours() {
	mdbok := len(dbdata.data.Info.MiddayBreak) == 0
	for i, d := range dbdata.w365.Hours {
		if d.FirstAfternoonHour {
			dbdata.data.Info.FirstAfternoonHour = i
		}
		if d.MiddayBreak {
			if mdbok {
				dbdata.data.Info.MiddayBreak = append(
					dbdata.data.Info.MiddayBreak, i)
			} else {
				fmt.Printf("*ERROR* MiddayBreak set in Info AND Hours")
			}
		}

		dbdata.data.Hours = append(dbdata.data.Hours, db.Hour{
			Id:    dbdata.nextId(d.Id),
			Type:  db.TypeHOUR,
			Tag:   d.Shortcut,
			Name:  d.Name,
			Start: d.Start,
			End:   d.End,
		})
	}
}

func DeMultipleSubjects(w365 *W365TopLevel) {
	/* Subjects -> Subject conversion */
	// First gather keys for all Subject nodes.
	subject2key := map[W365Ref]string{}
	for _, s := range w365.Subjects {
		subject2key[s.IdStr()] = s.Shortcut
	}
	cache := map[string]W365Ref{}
	// Now check all Courses and SubCourses for multiple subjects.
	n := 0
	for i, c := range w365.Courses {
		if c.Subject == "" {
			if len(c.Subjects) == 1 {
				w365.Courses[i].Subject = c.Subjects[0]
			} else if len(c.Subjects) > 1 {
				// Make a subject name
				sklist := []string{}
				for _, sid := range c.Subjects {
					sk, ok := subject2key[sid]
					if ok {
						sklist = append(sklist, sk)
					} else {
						fmt.Printf("*ERROR* Course %s:\n  Unknown Subject: %s\n",
							c.IdStr(), sid)
					}
				}
				skname := strings.Join(sklist, ",")
				sid, ok := cache[skname]
				if !ok {
					n++
					sk := fmt.Sprintf("X%02d", n)
					sid = W365Ref(fmt.Sprintf("Id_%s", sk))
					w365.Subjects = append(w365.Subjects, Subject{
						Id:       sid,
						Type:     TypeSUBJECT,
						Name:     skname,
						Shortcut: sk,
					})
					cache[skname] = sid
					subject2key[sid] = sk

				}
				w365.Courses[i].Subject = sid
			} else if len(c.Subjects) != 0 {
				fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n",
					c.IdStr())
			}
		}
	}
	for i, c := range w365.SubCourses {
		if c.Subject == "" {
			if len(c.Subjects) == 1 {
				w365.SubCourses[i].Subject = c.Subjects[0]
			} else if len(c.Subjects) > 1 {
				// Make a subject name
				sklist := []string{}
				for _, sid := range c.Subjects {
					sk, ok := subject2key[sid]
					if ok {
						sklist = append(sklist, sk)
					} else {
						fmt.Printf("*ERROR* Course %s:\n  Unknown Subject: %s\n",
							c.IdStr(), sid)
					}
				}
				skname := strings.Join(sklist, ",")
				sid, ok := cache[skname]
				if !ok {
					n++
					sk := fmt.Sprintf("X%02d", n)
					sid = W365Ref(fmt.Sprintf("Id_%s", sk))
					w365.Subjects = append(w365.Subjects, Subject{
						Id:       sid,
						Type:     TypeSUBJECT,
						Name:     skname,
						Shortcut: sk,
					})
					cache[skname] = sid
					subject2key[sid] = sk

				}
				w365.SubCourses[i].Subject = sid
			} else if len(c.Subjects) != 0 {
				fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n",
					c.IdStr())
			}
		}
	}
}
