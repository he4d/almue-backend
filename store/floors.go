package store

import (
	"github.com/he4d/almue/model"
)

func (d *datastore) GetFloor(floorID int64) (*model.Floor, error) {
	floor := &model.Floor{}
	//TODO: Add Shutters and Lightings
	err := d.QueryRow(floorFindIDStmt,
		floorID).Scan(&floor.ID, &floor.Created, &floor.Modified, &floor.Description)
	if err != nil {
		return nil, err
	}
	return floor, err
}

func (d *datastore) GetFloorList() ([]*model.Floor, error) {
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

func (d *datastore) CreateFloor(f *model.Floor) (int64, error) {
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

func (d *datastore) DeleteFloor(floorID int64) error {
	_, err := d.Exec(floorDeleteStmt, floorID)
	return err
}

func (d *datastore) UpdateFloor(f *model.Floor) error {
	_, err :=
		d.Exec(floorUpdateStmt, f.Description, f.ID)
	return err
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
