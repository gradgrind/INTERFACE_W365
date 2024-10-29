package w365tt

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func (dbdata *xData) readRooms() {
	for i := 0; i < len(dbdata.data.Rooms); i++ {
		n := &dbdata.data.Rooms[i]
		dbdata.elements[n.Id] = n
		if len(n.NotAvailable) == 0 {
			// Avoid a null value
			n.NotAvailable = []TimeSlot{}
		}
	}
}

func (dbdata *xData) readRoomGroups() {
	tags := map[string]bool{}
	tagless := []*RoomGroup{}
	for i := 0; i < len(dbdata.data.RoomGroups); i++ {
		n := &dbdata.data.RoomGroups[i]
		dbdata.elements[n.Id] = n

		n.Rooms = slices.DeleteFunc(n.Rooms, func(r Ref) bool {
			if rm, ok := dbdata.elements[r]; ok {
				if _, ok := rm.(*Room); ok {
					return false
				}
			}
			fmt.Printf("*ERROR* Unknown Room in RoomGroup %s:\n  %s\n",
				n.Tag, r)
			return true
		})

		if n.Tag == "" {
			tagless = append(tagless, n)
		} else {
			tags[n.Tag] = true
		}
	}
	for _, n := range tagless {
		rlist := []string{}
		for _, r := range n.Rooms {
			rlist = append(rlist, dbdata.elements[r].(*Room).Tag)
		}
		tag := fmt.Sprintf("{%s}", strings.Join(rlist, ","))
		i := 1
		if tags[tag] {
			for {
				ti := tag + strconv.Itoa(i)
				if !tags[ti] {
					tag = ti
					tags[ti] = true
					break
				}
				i++
			}
		}
		n.Tag = tag
	}
}

func (dbdata *xData) readRoomChoiceGroups() {
	tags := map[string]bool{}
	tagless := []*RoomChoiceGroup{}
	for i := 0; i < len(dbdata.data.RoomChoiceGroups); i++ {
		n := &dbdata.data.RoomChoiceGroups[i]
		dbdata.elements[n.Id] = n

		n.Rooms = slices.DeleteFunc(n.Rooms, func(r Ref) bool {
			if rm, ok := dbdata.elements[r]; ok {
				if _, ok := rm.(*Room); ok {
					return false
				}
			}
			fmt.Printf("*ERROR* Unknown Room in RoomChoiceGroup %s:\n  %s\n",
				n.Tag, r)
			return true
		})

		if n.Tag == "" {
			tagless = append(tagless, n)
		} else {
			tags[n.Tag] = true
		}
	}
	for _, n := range tagless {
		rlist := []string{}
		for _, r := range n.Rooms {
			rlist = append(rlist, dbdata.elements[r].(*Room).Tag)
		}
		tag := fmt.Sprintf("[%s]", strings.Join(rlist, ","))
		i := 1
		if tags[tag] {
			for {
				ti := tag + strconv.Itoa(i)
				if !tags[ti] {
					tag = ti
					tags[ti] = true
					break
				}
				i++
			}
		}
		n.Tag = tag
	}
}
