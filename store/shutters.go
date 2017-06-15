package store

import (
	"errors"

	"github.com/he4d/almue/model"
)

func (d *datastore) GetShutterByFloor(floorID, shutterID int64) (*model.Shutter, error) {
	return nil, errors.New("not implemented")
}

func (d *datastore) GetShutterListOfFloor(floorID int64) ([]*model.Shutter, error) {
	return nil, errors.New("not implemented")
}

func (d *datastore) GetAllShutters() ([]*model.Shutter, error) {
	return nil, errors.New("not implemented")
}

func (d *datastore) CreateShutter(*model.Shutter) error {
	return errors.New("not implemented")
}

func (d *datastore) DeleteShutter(*model.Shutter) error {
	return errors.New("not implemented")
}

func (d *datastore) UpdateShutter(*model.Shutter) error {
	return errors.New("not implemented")
}

var shutterFindID = `
SELECT * FROM shutters WHERE id = ? AND floor_id = ?
`

var shuttersOfFloor = `
SELECT * FROM shutters WHERE floor_id = ?
`

var shuttersFindAll = `
SELECT * FROM shutters
`

var shutterCreate = `
INSERT INTO shutters(
description,
open_pin,
close_pin,
complete_way_in_seconds,
timer_enabled,
open_time,
close_time,
emergency_enabled,
device_status,
disabled,
floor_id
) 
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

var shutterUpdate = `
UPDATE shutters SET 
description = ?,
open_pin = ?,
close_pin = ?,
complete_way_in_seconds = ?,
timer_enabled = ?,
open_time = ?,
close_time = ?,
emergency_enabled = ?,
device_status = ?,
disabled = ? 
WHERE id = ?
`

var shutterDelete = `
DELETE FROM shutters where id = ? and floor_id = ?
`
