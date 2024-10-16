package w365_tt

import "encoding/xml"

// The structures used for reading a timetable-source file exported by W365.

type Day struct {
	Id           string
	ListPosition float32
	Name         string
	Shortcut     string
}

type W365TT struct {
	XMLName xml.Name `xml:"File"`
	Days    []Day
}
