package w365tt

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

func (dbp *DbTopLevel) readRooms() {
	for i := 0; i < len(dbp.Rooms); i++ {
		n := &dbp.Rooms[i]
		if len(n.NotAvailable) == 0 {
			// Avoid a null value
			n.NotAvailable = []TimeSlot{}
		}
	}
}

func (dbp *DbTopLevel) readRoomGroups() {
	tags := map[string]bool{}
	tagless := []*RoomGroup{}
	for i := 0; i < len(dbp.RoomGroups); i++ {
		n := &dbp.RoomGroups[i]

		n.Rooms = slices.DeleteFunc(n.Rooms, func(r Ref) bool {
			if rm, ok := dbp.Elements[r]; ok {
				if _, ok := rm.(*Room); ok {
					return false
				}
			}
			log.Printf("*ERROR* Unknown Room in RoomGroup %s:\n  %s\n",
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
			rlist = append(rlist, dbp.Elements[r].(*Room).Tag)
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

func (dbp *DbTopLevel) readRoomChoiceGroups() {
	tags := map[string]bool{}
	tagless := []*RoomChoiceGroup{}
	for i := 0; i < len(dbp.RoomChoiceGroups); i++ {
		n := &dbp.RoomChoiceGroups[i]

		n.Rooms = slices.DeleteFunc(n.Rooms, func(r Ref) bool {
			if rm, ok := dbp.Elements[r]; ok {
				if _, ok := rm.(*Room); ok {
					return false
				}
			}
			log.Printf("*ERROR* Unknown Room in RoomChoiceGroup %s:\n  %s\n",
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
			rlist = append(rlist, dbp.Elements[r].(*Room).Tag)
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
