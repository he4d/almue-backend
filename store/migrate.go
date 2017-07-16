package store

import "database/sql"

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-floors",
		stmt: createTableFloors,
	},
	{
		name: "create-table-shutters",
		stmt: createTableShutters,
	},
	{
		name: "create-table-lightings",
		stmt: createTableLightings,
	},
	{
		name: "create-update-trigger-floors",
		stmt: createUpdateTriggerFloors,
	},
	{
		name: "create-update-trigger-shutters",
		stmt: createUpdateTriggerShutters,
	},
	{
		name: "create-update-trigger-lightings",
		stmt: createUpdateTriggerLightings,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
name VARCHAR(255),
UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

var createTableFloors = `
CREATE TABLE IF NOT EXISTS floors (
id integer primary key,
created datetime NOT NULL DEFAULT current_timestamp,
modified datetime NOT NULL DEFAULT current_timestamp,
description varchar(255) NOT NULL UNIQUE 
)
`

//TODO: The pins must be unique merged (open_pin and close_pin) TABLE FOR PINS!!
var createTableShutters = `
CREATE TABLE IF NOT EXISTS shutters (
id integer primary key,
created datetime NOT NULL DEFAULT current_timestamp,
modified datetime NOT NULL DEFAULT current_timestamp,
description varchar(255),
open_pin integer NOT NULL UNIQUE,
close_pin integer NOT NULL UNIQUE,
complete_way_in_seconds integer NOT NULL,
opening_in_prc integer DEFAULT 0,
jobs_enabled bool,
open_time datetime DEFAULT current_timestamp,
close_time datetime DEFAULT current_timestamp,
emergency_enabled bool,
device_status varchar(10),
disabled bool,
floor_id integer NOT NULL REFERENCES floors(id) ON DELETE CASCADE ON UPDATE CASCADE 
)
`

//TODO: The pins must be unique merged with shutters (open_pin and close_pin)
var createTableLightings = `
CREATE TABLE IF NOT EXISTS lightings (
id integer primary key,
created datetime NOT NULL DEFAULT current_timestamp,
modified datetime NOT NULL DEFAULT current_timestamp,
description varchar(255),
switch_pin integer NOT NULL UNIQUE,
jobs_enabled bool,
on_time datetime DEFAULT current_timestamp,
off_time datetime DEFAULT current_timestamp,
emergency_enabled bool,
device_status varchar(10),
disabled bool,
floor_id integer NOT NULL REFERENCES floors(id) ON DELETE CASCADE ON UPDATE CASCADE
)
`

var createUpdateTriggerFloors = `
CREATE TRIGGER IF NOT EXISTS 
update_floor AFTER UPDATE ON floors FOR EACH ROW BEGIN UPDATE floors 
SET modified = current_timestamp WHERE ID = OLD.ID; END;
`

var createUpdateTriggerShutters = `
CREATE TRIGGER IF NOT EXISTS 
update_shutter AFTER UPDATE ON shutters FOR EACH ROW BEGIN UPDATE shutters 
SET modified = current_timestamp WHERE ID = OLD.ID; END;
`

var createUpdateTriggerLightings = `
CREATE TRIGGER IF NOT EXISTS 
update_lighting AFTER UPDATE ON lightings FOR EACH ROW BEGIN UPDATE lightings 
SET modified = current_timestamp WHERE ID = OLD.ID; END;
`
