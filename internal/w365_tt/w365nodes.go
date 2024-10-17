package w365_tt

import "encoding/xml"

// The structures used for reading a timetable-source file exported by W365.

type Day struct {
	Id           string  `xml:",attr"`
	ListPosition float32 `xml:",attr"`
	Name         string  `xml:",attr"`
	Shortcut     string  `xml:",attr"`
}

type W365TT struct {
	XMLName xml.Name `xml:"File"`
	Days    []Day    `xml:"Day"`
}
