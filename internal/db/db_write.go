package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
I need a list of database nodes, presumably containing their own DB-key.
The question is how they get built ...
*/

type Node map[string]interface{}

func (n Node) GetId() int {
	return n["Id"].(int)
}

func get_days(nlist *[]Node) {
	*nlist = append(*nlist, Node{
		"Id":   1,
		"NAME": "Montag",
		"TAG":  "Mo",
		"TYPE": "DAYS",
		"X":    0,
	})
	*nlist = append(*nlist, Node{
		"Id":   2,
		"NAME": "Dienstag",
		"TAG":  "Di",
		"TYPE": "DAYS",
		"X":    1,
	})
}

func test1(dbpath string) {
	/*	fmt.Printf("TIME (22:7:12): %s\n", get_time("22:7:12"))
		fmt.Printf("TIME (6:07): %s\n", get_time("6:07"))
		fmt.Printf("TIME (24:07): %s\n", get_time("24:07"))
		fmt.Printf("TIME (-2:07): %s\n", get_time("-2:07"))
	*/
	nodelist := []Node{}
	get_days(&nodelist)
	dbwrite(dbpath, nodelist)
}

func dbwrite(dbpath string, nodelist []Node) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
DROP TABLE IF EXISTS NODES;
CREATE TABLE NODES(
	Id INTEGER PRIMARY KEY,
	DATA TEXT NOT NULL
);
`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = "INSERT INTO NODES(Id, DATA) values(?,?)"
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// The primary key will correspond to the node indexes.
	for _, node := range nodelist {
		j, err := json.Marshal(node)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		fmt.Printf(" ??? %+v\n", node)
		_, err = tx.Exec(query, node.GetId(), string(j))
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
