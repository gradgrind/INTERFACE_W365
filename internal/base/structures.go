package base

import (
	"encoding/json"
	"fmt"
	"log"
)

// A "top-level" JSON object, a database DATA value
type Record map[string]interface{}

func (n Record) GetId() int {
	return n["Id"].(int) // If the field is missing this will panic.
}

func (n Record) GetType() string {
	return n["Type"].(string) // If the field is missing this will panic.
}

func (n Record) GetX() int {
	x, ok := n["X"].(int)
	if ok {
		return x
	}
	return -1
}

func SaveJSON(records []Record, jsonpath string) {
	j, err := json.MarshalIndent(records, "++", "  ")
	//j, err := json.Marshal(records)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" ??? %+v\n", string(j))
}

type DBData struct {
	Records []Record // The first entry should be invalid (e.g. empty).
	//	Id2Index map[int]int      // Map Id to Record index
	Tables map[string][]int // Map record type to list of Record indexes
}

func NewDBData() DBData {
	var d DBData
	d.Records = []Record{{}}
	d.Tables = make(map[string][]int)
	return d
}

func (dbdata *DBData) SetInfo(key string, val interface{}) {
	dbdata.Records[0][key] = val
}

func (dbdata *DBData) AddInfo(key string, val interface{}) {
	l := dbdata.Records[0][key]
	if l == nil {
		dbdata.Records[0][key] = []interface{}{val}
	} else {
		dbdata.Records[0][key] = append(l.([]interface{}), val)
	}
}

/*
func (dbdata *DBData) Id2Tag(id int) string {
	return dbdata.Records[dbdata.Id2Index[id]]["Tag"].(string)
}
*/

func (dbdata *DBData) AddRecord(r Record) int {
	i := len(dbdata.Records)
	dbdata.Records = append(dbdata.Records, r)
	r["Id"] = i
	t := r.GetType()
	dbdata.Tables[t] = append(dbdata.Tables[t], i)
	return i
}

const (
	RecordType_DAY     string = "DAY"
	RecordType_HOUR    string = "HOUR"
	RecordType_SUBJECT string = "SUBJECT"
	RecordType_TEACHER string = "TEACHER"
	RecordType_ROOM    string = "ROOM"
	RecordType_CLASS   string = "CLASS"
	RecordType_GROUP   string = "GROUP"
)
