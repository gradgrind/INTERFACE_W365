package fet

import (
	"encoding/xml"
)

// TODO: At present this is used for adding activity-tags to activities
// for these subjects (see courses.go). This is just a temporary hack ...
var tagged_subjects = []string{"Eu", "Sp"}

type fetSubject struct {
	XMLName   xml.Name `xml:"Subject"`
	Name      string
	Long_Name string
	Comments  string
}

type fetSubjectsList struct {
	XMLName xml.Name `xml:"Subjects_List"`
	Subject []fetSubject
}

func getSubjects(fetinfo *fetInfo) {
	items := []fetSubject{}
	for _, n := range fetinfo.db.Subjects {
		items = append(items, fetSubject{
			Name:      n.Tag,
			Long_Name: n.Name,
		})
	}
	fetinfo.fetdata.Subjects_List = fetSubjectsList{
		Subject: items,
	}
}
