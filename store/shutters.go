package store

import (
	"fmt"

	"github.com/he4d/almue/model"
)

func (d *Datastore) GetShutter(shutterID int64) (*model.Shutter, error) {
	s := new(model.Shutter)

	err := d.QueryRow(shutterByIDStmt, shutterID).Scan(
		&s.ID, &s.Created, &s.Modified, &s.Description,
		&s.OpenPin, &s.ClosePin, &s.CompleteWayInSeconds,
		&s.OpeningInPrc, &s.JobsEnabled, &s.OpenTime, &s.CloseTime,
		&s.EmergencyEnabled, &s.DeviceStatus, &s.Disabled,
		&s.FloorID)

	if err != nil {
		return nil, err
	}
	return s, err
}

func (d *Datastore) GetShutterListOfFloor(floorID int64) ([]*model.Shutter, error) {
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
			&s.OpeningInPrc, &s.JobsEnabled, &s.OpenTime, &s.CloseTime,
			&s.EmergencyEnabled, &s.DeviceStatus, &s.Disabled,
			&s.FloorID); err != nil {
			return nil, err
		}
		shutters = append(shutters, s)
	}

	return shutters, err
}

func (d *Datastore) GetShutterList() ([]*model.Shutter, error) {
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
			&s.OpeningInPrc, &s.JobsEnabled, &s.OpenTime, &s.CloseTime,
			&s.EmergencyEnabled, &s.DeviceStatus, &s.Disabled,
			&s.FloorID); err != nil {
			return nil, err
		}
		shutters = append(shutters, s)
	}

	return shutters, err
}

func (d *Datastore) CreateShutter(s *model.Shutter) (int64, error) {
	res, err := d.Exec(
		shutterCreateStmt,
		s.Description, s.OpenPin, s.ClosePin, s.CompleteWayInSeconds,
		s.JobsEnabled, s.OpenTime.UTC(), s.CloseTime.UTC(), s.EmergencyEnabled,
		"stopped", s.Disabled, s.FloorID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (d *Datastore) DeleteShutter(shutterID int64) error {
	res, err := d.Exec(shutterDeleteStmt, shutterID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("Shutter with id %d didnt exist", shutterID)
	}
	return err
}

func (d *Datastore) UpdateShutter(s *model.Shutter) error {
	_, err :=
		d.Exec(
			shutterUpdateStmt,
			s.Description, s.OpenPin, s.ClosePin, s.CompleteWayInSeconds,
			s.JobsEnabled, s.OpenTime.UTC(), s.CloseTime.UTC(), s.EmergencyEnabled,
			s.DeviceStatus, s.Disabled, s.FloorID, s.ID)
	return err
}

func (d *Datastore) UpdateShutterState(shutterID int64, newState string) error {
	_, err :=
		d.Exec(shutterStateUpdateStmt, newState, shutterID)
	return err
}

func (d *Datastore) UpdateShutterOpening(shutterID int64, openingInPrc float64) error {
	_, err :=
		d.Exec(shutterOpeningInPrcUpdateStmt, openingInPrc, shutterID)
	return err
}

var shutterByIDStmt = `
SELECT * FROM shutters WHERE id = ?
`

var shuttersOfFloorStmt = `
SELECT * FROM shutters WHERE floor_id = ?
`

var shuttersFindAllStmt = `
SELECT * FROM shutters
`

var shutterStateUpdateStmt = `
UPDATE shutters SET
device_status = ? 
WHERE id = ?
`

var shutterOpeningInPrcUpdateStmt = `
UPDATE shutters SET
opening_in_prc = ?
WHERE id = ?
`

var shutterCreateStmt = `
INSERT INTO shutters(
description,
open_pin,
close_pin,
complete_way_in_seconds,
jobs_enabled,
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
jobs_enabled = ?,
open_time = ?,
close_time = ?,
emergency_enabled = ?,
device_status = ?,
disabled = ?,
floor_id = ? 
WHERE id = ?
`

var shutterDeleteStmt = `
DELETE FROM shutters WHERE id = ?
`
