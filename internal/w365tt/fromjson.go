package w365tt

import (
	"encoding/json"
	"fmt"
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
	return v
}

// TODO: This should probably be in the fet package!
type FETData struct {
}

// Subjects -> Subject conversion
func LoadJSON(jsonpath string) FETData {
	toplevel := ReadJSON(jsonpath)
	fetdata := FETData{}
	return fetdata
}

func DeMultipleSubjects(w365 *W365TopLevel) {
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
}
