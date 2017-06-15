package store

import (
	"errors"

	"github.com/he4d/almue/model"
)

func (d *datastore) GetLightingByFloor(floorID, lightingID int64) (*model.Lighting, error) {
	return nil, errors.New("not implemented")
}

func (d *datastore) GetLightingListOfFloor(floorID int64) ([]*model.Lighting, error) {
	return nil, errors.New("not implemented")
}

func (d *datastore) GetAllLightings() ([]*model.Lighting, error) {
	return nil, errors.New("not implemented")
}

func (d *datastore) CreateLighting(*model.Lighting) error {
	return errors.New("not implemented")
}

func (d *datastore) DeleteLighting(*model.Lighting) error {
	return errors.New("not implemented")
}

func (d *datastore) UpdateLighting(*model.Lighting) error {
	return errors.New("not implemented")
}

var lightingByID = `
SELECT * FROM lightings WHERE id = ? AND floor_id = ?
`
var lightingsOfFloor = `
SELECT * FROM lightings WHERE floor_id = ?
`

var lightingsFindAll = `
SELECT * FROM lightings
`

var lightingCreate = `
INSERT INTO lightings(
description,
switch_pin,
timer_enabled,
on_time,
off_time,
emergency_enabled,
device_status,
disabled,
floor_id
) 
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
`

var lightingUpdate = `
UPDATE lightings SET 
description = ?,
switch_pin = ?,
timer_enabled = ?,
on_time = ?,
off_time = ?,
emergency_enabled = ?,
device_status = ?,
disabled = ?  
WHERE id = ?
`

var lightingDelete = `
DELETE FROM lightings WHERE id = ?
`
