package store

import (
	"github.com/he4d/almue/model"
)

func (d *datastore) GetShutterByFloor(shutterID, floorID int64) (*model.Shutter, error) {
	s := new(model.Shutter)

	err := d.QueryRow(shutterFindIDStmt, shutterID, floorID).Scan(
		&s.ID, &s.Created, &s.Modified, &s.Description,
		&s.OpenPin, &s.ClosePin, &s.CompleteWayInSeconds,
		&s.TimerEnabled, &s.OpenTime, &s.CloseTime,
		&s.EmergencyEnabled, &s.DeviceStatus, &s.Disabled,
		&s.FloorID)

	if err != nil {
		return nil, err
	}
	return s, err
}

func (d *datastore) GetShutterListOfFloor(floorID int64) ([]*model.Shutter, error) {
	rows, err := d.Query(shuttersOfFloorStmt, floorID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shutters := []*model.Shutter{}

	for rows.Next() {
		s := new(model.Shutter)
		if err := rows.Scan(
			&s.ID, &s.Created, &s.Modified, &s.Description,
			&s.OpenPin, &s.ClosePin, &s.CompleteWayInSeconds,
			&s.TimerEnabled, &s.OpenTime, &s.CloseTime,
			&s.EmergencyEnabled, &s.DeviceStatus, &s.Disabled,
			&s.FloorID); err != nil {
			return nil, err
		}
		shutters = append(shutters, s)
	}

	return shutters, err
}

func (d *datastore) GetAllShutters() ([]*model.Shutter, error) {
	rows, err := d.Query(shuttersFindAllStmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shutters := []*model.Shutter{}

	for rows.Next() {
		s := new(model.Shutter)
		if err := rows.Scan(
			&s.ID, &s.Created, &s.Modified, &s.Description,
			&s.OpenPin, &s.ClosePin, &s.CompleteWayInSeconds,
			&s.TimerEnabled, &s.OpenTime, &s.CloseTime,
			&s.EmergencyEnabled, &s.DeviceStatus, &s.Disabled,
			&s.FloorID); err != nil {
			return nil, err
		}
		shutters = append(shutters, s)
	}

	return shutters, err
}

func (d *datastore) CreateShutter(s *model.Shutter) (int64, error) {
	res, err := d.Exec(
		shutterCreateStmt,
		s.Description, s.OpenPin, s.ClosePin, s.CompleteWayInSeconds,
		s.TimerEnabled, s.OpenTime, s.CloseTime, s.EmergencyEnabled,
		s.DeviceStatus, s.Disabled, s.FloorID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (d *datastore) DeleteShutter(shutterID int64) error {
	_, err := d.Exec(shutterDeleteStmt, shutterID)
	return err
}

func (d *datastore) UpdateShutter(s *model.Shutter) error {
	_, err :=
		d.Exec(
			shutterUpdateStmt,
			s.Description, s.OpenPin, s.ClosePin, s.CompleteWayInSeconds,
			s.TimerEnabled, s.OpenTime, s.CloseTime, s.EmergencyEnabled,
			s.DeviceStatus, s.Disabled, s.ID)
	return err
}

var shutterFindIDStmt = `
SELECT * FROM shutters WHERE id = ? AND floor_id = ?
`

var shuttersOfFloorStmt = `
SELECT * FROM shutters WHERE floor_id = ?
`

var shuttersFindAllStmt = `
SELECT * FROM shutters
`

var shutterCreateStmt = `
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

var shutterUpdateStmt = `
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

var shutterDeleteStmt = `
DELETE FROM shutters WHERE id = ?
`
