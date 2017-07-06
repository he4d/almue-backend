package store

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetEmptyFloorsTable(t *testing.T) {
	clearTable()

	floors, err := store.GetFloorList()
	if err != nil {
		t.Error(err)
	}
	if len(floors) != 0 {
		t.Error("Floors table not empty")
	}
}

func TestGetSingleFloor(t *testing.T) {
	clearTable()

	descr := "obergeschoss"
	res, err := store.Exec("INSERT INTO floors(description) VALUES(?)", descr)
	if err != nil {
		t.Errorf("Could not create the init floor: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Errorf("Did not get the last inserted id on the response: %v", err)
	}
	floor, err := store.GetFloor(id)
	if err != nil {
		t.Errorf("Could not get the floor for testing: %v", err)
	}
	if floor == nil {
		t.Error("Got no error before but floor is nil")
	}
	if *floor.Description != descr {
		t.Errorf("Got the floor with a wrong description: %s", *floor.Description)
	}
}

func TestGetFloorList(t *testing.T) {
	clearTable()
	const amountOfFloors = 5
	const floorDescr = "testfloor"
	for i := 0; i < amountOfFloors; i++ {
		descr := fmt.Sprintf("%s%d", floorDescr, i)
		_, err := store.Exec("INSERT INTO floors(description) VALUES(?)", descr)
		if err != nil {
			t.Errorf("Could not create the init floor: %v", err)
		}
	}

	floors, err := store.GetFloorList()
	if err != nil {
		t.Errorf("Could not get the floor list %v", err)
	}
	if floors == nil {
		t.Error("Floors is nil but no error occured before")
	}
	numOfFloors := len(floors)
	if numOfFloors != amountOfFloors {
		t.Errorf("%d floors created but got %d", amountOfFloors, numOfFloors)
	}
	for idx, floor := range floors {
		descr := fmt.Sprintf("%s%d", floorDescr, idx)
		if *floor.Description != descr {
			t.Errorf("Created floor with description %s but got %s", descr, *floor.Description)
		}
	}
}

func clearTable() {
	_, err := store.Exec("DELETE FROM floors")
	if err != nil {
		log.Fatalf("Could not clear the table: %v", err)
	}
}
