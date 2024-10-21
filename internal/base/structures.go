package base

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

/*
func (dbdata *DBData) Id2Tag(id int) string {
	return dbdata.Records[dbdata.Id2Index[id]]["Tag"].(string)
}
*/

func (dbdata *DBData) addRecord(r Record) int {
	i := len(dbdata.Records)
	dbdata.Records = append(dbdata.Records, r)
	r["Id"] = i
	t := r.GetType()
	dbdata.Tables[t] = append(dbdata.Tables[t], i)
	return i
}
