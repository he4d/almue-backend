package store

import (
	"fmt"

	"github.com/he4d/almue/model"
)

func (d *Datastore) GetFloor(floorID int64) (*model.Floor, error) {
	floor := &model.Floor{}
	err := d.QueryRow(floorFindIDStmt,
		floorID).Scan(&floor.ID, &floor.Created, &floor.Modified, &floor.Description)
	if err != nil {
		return nil, err
	}
	return floor, err
}

func (d *Datastore) GetFloorList() ([]*model.Floor, error) {
	rows, err := d.Query(floorsFindAllStmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	floors := []*model.Floor{}

	for rows.Next() {
		var f model.Floor
		if err := rows.Scan(&f.ID, &f.Created, &f.Modified, &f.Description); err != nil {
			return nil, err
		}
		floors = append(floors, &f)
	}
	return floors, err
}

func (d *Datastore) CreateFloor(f *model.Floor) (int64, error) {
	res, err := d.Exec(
		floorCreateStmt,
		f.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (d *Datastore) DeleteFloor(floorID int64) error {
	res, err := d.Exec(floorDeleteStmt, floorID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("Floor with id %d didnt exist", floorID)
	}
	return err
}

func (d *Datastore) UpdateFloor(f *model.Floor) error {
	_, err :=
		d.Exec(floorUpdateStmt, f.Description, f.ID)
	return err
}

func (d *Datastore) NumShuttersOfFloor(floorID int64) (int, error) {
	retval := 0
	err := d.QueryRow(numShuttersOfFloorStmt,
		floorID).Scan(&retval)
	if err != nil {
		return 0, err
	}
	return retval, err
}

func (d *Datastore) NumLightingsOfFloor(floorID int64) (int, error) {
	retval := 0
	err := d.QueryRow(numLightingsfFloorStmt,
		floorID).Scan(&retval)
	if err != nil {
		return 0, err
	}
	return retval, err
}

var floorFindIDStmt = `
SELECT * FROM floors WHERE id = ?
`

var floorsFindAllStmt = `
SELECT * FROM floors
`

var floorCreateStmt = `
INSERT INTO floors(description) VALUES(?)
`

var floorUpdateStmt = `
UPDATE floors SET description = ? WHERE id = ?
`

var floorDeleteStmt = `
DELETE FROM floors WHERE id = ?
`

var numShuttersOfFloorStmt = `
SELECT COUNT(id) FROM shutters where floor_id = ?
`

var numLightingsfFloorStmt = `
SELECT COUNT(id) FROM lightings where floor_id = ?
`
