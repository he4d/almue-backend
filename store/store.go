package store

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"io/ioutil"

	"github.com/he4d/simplejack"
	sqlite3 "github.com/mattn/go-sqlite3"
)

// Datastore contains all necessary objects for store handling
type Datastore struct {
	*sql.DB
	logger *simplejack.Logger
}

// New returns a new datastore that is completely initialized
func New(path string, logger *simplejack.Logger) (*Datastore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := setupDatabase(db); err != nil {
		return nil, err
	}
	return &Datastore{DB: db, logger: logger}, nil
}

func setupDatabase(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	return Migrate(db)
}

// GetBackup creates a database backup and returns it as a byte array
func (d *Datastore) GetBackup() ([]byte, error) {
	var driverName = fmt.Sprintf("sqlite3_backup_%v", time.Now().UnixNano())
	tmpFile, err := ioutil.TempFile("", "tmpDb")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	// The driver's connection will be needed in order to perform the backup.
	driverConns := []*sqlite3.SQLiteConn{}

	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			driverConns = append(driverConns, conn)
			return nil
		},
	})

	// Connect to the source database.
	srcFileName := "almue.db"
	srcDb, err := sql.Open(driverName, srcFileName)
	if err != nil {
		return nil, err
	}
	defer srcDb.Close()
	err = srcDb.Ping()
	if err != nil {
		return nil, err
	}

	// Connect to the destination database.
	destDb, err := sql.Open(driverName, tmpFile.Name())
	if err != nil {
		return nil, err
	}
	defer destDb.Close()
	err = destDb.Ping()
	if err != nil {
		return nil, err
	}

	if len(driverConns) != 2 {
		return nil, fmt.Errorf("Expected 2 driver connections, but found %v", len(driverConns))
	}
	srcDbDriverConn := driverConns[0]
	if srcDbDriverConn == nil {
		return nil, err
	}
	destDbDriverConn := driverConns[1]
	if destDbDriverConn == nil {
		return nil, err
	}

	backup, err := destDbDriverConn.Backup("main", srcDbDriverConn, "main")
	if err != nil {
		return nil, err
	}

	isDone, err := backup.Step(-1)
	if err != nil {
		return nil, err
	}
	if !isDone {
		return nil, fmt.Errorf("Backup is unexpectedly not done")
	}

	err = backup.Finish()
	if err != nil {
		return nil, err
	}

	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
