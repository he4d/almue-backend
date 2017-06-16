package store

import (
	"github.com/he4d/almue/model"
)

func (d *datastore) GetLightingByFloor(lightingID, floorID int64) (*model.Lighting, error) {
	l := new(model.Lighting)

	err := d.QueryRow(lightingByIDStmt, lightingID, floorID).Scan(
		&l.ID, &l.Created, &l.Modified, &l.Description,
		&l.SwitchPin, &l.TimerEnabled, &l.OnTime, &l.OffTime,
		&l.EmergencyEnabled, &l.DeviceStatus, &l.Disabled,
		&l.FloorID)

	if err != nil {
		return nil, err
	}
	return l, err
}

func (d *datastore) GetLightingListOfFloor(floorID int64) ([]*model.Lighting, error) {
	rows, err := d.Query(lightingsOfFloorStmt, floorID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	lightings := []*model.Lighting{}

	for rows.Next() {
		l := new(model.Lighting)
		if err := rows.Scan(
			&l.ID, &l.Created, &l.Modified, &l.Description,
			&l.SwitchPin, &l.TimerEnabled, &l.OnTime, &l.OffTime,
			&l.EmergencyEnabled, &l.DeviceStatus, &l.Disabled,
			&l.FloorID); err != nil {
			return nil, err
		}
		lightings = append(lightings, l)
	}

	return lightings, err
}

func (d *datastore) GetAllLightings() ([]*model.Lighting, error) {
	rows, err := d.Query(lightingsFindAllStmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	lightings := []*model.Lighting{}

	for rows.Next() {
		l := new(model.Lighting)
		if err := rows.Scan(
			&l.ID, &l.Created, &l.Modified, &l.Description,
			&l.SwitchPin, &l.TimerEnabled, &l.OnTime, &l.OffTime,
			&l.EmergencyEnabled, &l.DeviceStatus, &l.Disabled,
			&l.FloorID); err != nil {
			return nil, err
		}
		lightings = append(lightings, l)
	}

	return lightings, err
}

func (d *datastore) CreateLighting(l *model.Lighting) (int64, error) {
	res, err := d.Exec(
		lightingCreateStmt,
		l.Description, l.SwitchPin, l.TimerEnabled,
		l.OnTime, l.OffTime, l.EmergencyEnabled,
		l.DeviceStatus, l.Disabled, l.FloorID)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (d *datastore) DeleteLighting(lightingID int64) error {
	_, err := d.Exec(lightingDeleteStmt, lightingID)
	return err
}

func (d *datastore) UpdateLighting(l *model.Lighting) error {
	_, err :=
		d.Exec(
			lightingUpdateStmt,
			l.Description, l.SwitchPin,
			l.TimerEnabled, l.OnTime, l.OffTime, l.EmergencyEnabled,
			l.DeviceStatus, l.Disabled, l.ID)
	return err
}

var lightingByIDStmt = `
SELECT * FROM lightings WHERE id = ? AND floor_id = ?
`
var lightingsOfFloorStmt = `
SELECT * FROM lightings WHERE floor_id = ?
`

var lightingsFindAllStmt = `
SELECT * FROM lightings
`

var lightingCreateStmt = `
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

var lightingUpdateStmt = `
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

var lightingDeleteStmt = `
DELETE FROM lightings WHERE id = ?
`
