package base

import (
	"fmt"
	"testing"
)

func TestStructures(t *testing.T) {
	db_data := NewDBData()
	db_data.addRecord(Record{"Type": "TEST1"})
	db_data.addRecord(Record{"Type": "TEST2", "X": 3})

	fmt.Printf("STATE: %+v\n", db_data)
	// fmt.Printf("%s (%d): %d\n", r1["Type"], r1.GetId(), r1.GetX())
	// fmt.Printf("%s (%d): %d\n", r2["Type"], r2.GetId(), r2.GetX())
}
