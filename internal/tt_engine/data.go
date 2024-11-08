package tt_engine

import "gradgrind/INTERFACE_W365/internal/core"

type TtActivity struct {
	Length         int
	Resources      []Resource
	Lessons        []core.Lesson
	PermittedSlots []int
}

type Resource int

type TtData struct {
	ResourceMap  map[core.Ref]int
	ResourceList []any
}

// ???
func (ttdata *TtData) getElement(ref core.Ref) any {
	rmi, ok := ttdata.ResourceMap[ref]
	if !ok {
		core.Error.Fatalf("Unknown Element: %s\n", ref)
	}
	return ttdata.ResourceList[rmi]
}

func initData(db *core.DbTopLevel) *TtData {
	ttData := &TtData{
		ResourceMap:  map[core.Ref]int{},
		ResourceList: []any{nil}, // index 0 is invalid
	}

	// Teachers
	for i := 0; i < len(db.Teachers); i++ {
		n := &db.Teachers[i]
		ttData.ResourceMap[n.Id] = len(ttData.ResourceList)
		ttData.ResourceList = append(ttData.ResourceList, n)
	}

	// Real Rooms
	for i := 0; i < len(db.Rooms); i++ {
		n := &db.Rooms[i]
		ttData.ResourceMap[n.Id] = len(ttData.ResourceList)
		ttData.ResourceList = append(ttData.ResourceList, n)
	}

	// Classes and Groups
	// This is much more complicated as Atomic Groups are needed here.
	// First generate these.

	return ttData
}
