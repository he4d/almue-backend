package store

import (
	"errors"

	"github.com/he4d/almue/model"
)

func (d *datastore) GetFloor(floorID int64) (*model.Floor, error) {
	floor := &model.Floor{}
	//TODO: Add Shutters and Lightings
	err := d.QueryRow(floorFindID,
		floorID).Scan(&floor.ID, &floor.Created, &floor.Modified, &floor.Description)
	if err != nil {
		return nil, err
	}
	return floor, nil
}

func (d *datastore) GetFloorList() ([]*model.Floor, error) {
	rows, err := d.Query(floorsFindAll)

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
	return floors, nil
}

func (d *datastore) CreateFloor(*model.Floor) error {
	return errors.New("not implemented")
}

func (d *datastore) DeleteFloor(*model.Floor) error {
	return errors.New("not implemented")
}

func (d *datastore) UpdateFloor(*model.Floor) error {
	return errors.New("not implemented")
}

var floorFindID = `
SELECT * FROM floors WHERE id = ?
`

var floorsFindAll = `
SELECT * FROM floors
`

var floorCreate = `
INSERT INTO floors(description) VALUES(?)
`

var floorUpdate = `
UPDATE floors SET description = ? WHERE id = ?
`

var floorDelete = `
DELETE FROM floors WHERE id = ?
`
