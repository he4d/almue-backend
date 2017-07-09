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

//TODO: Divide and Conquer
func (d *Datastore) GetBackup() ([]byte, error) {
	var driverName = fmt.Sprintf("sqlite3_backup_%v", time.Now().UnixNano())
	tmpFile, err := ioutil.TempFile("", "tmpDb")
	if err != nil {
		return nil, fmt.Errorf("Could not create the tempFile for creating the backup: %v", err)
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
		return nil, fmt.Errorf("Failed to open the source db: %v", err)
	}
	defer srcDb.Close()
	err = srcDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the source database: %v", err)
	}

	// Connect to the destination database.
	destDb, err := sql.Open(driverName, tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("Failed to open the destination database: %v", err)
	}
	defer destDb.Close()
	err = destDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the destination database: %v", err)
	}

	if len(driverConns) != 2 {
		return nil, fmt.Errorf("Expected 2 driver connections, but found %v", len(driverConns))
	}
	srcDbDriverConn := driverConns[0]
	if srcDbDriverConn == nil {
		return nil, fmt.Errorf("The source database driver connection is nil")
	}
	destDbDriverConn := driverConns[1]
	if destDbDriverConn == nil {
		return nil, fmt.Errorf("The destination database driver connection is nil")
	}

	backup, err := destDbDriverConn.Backup("main", srcDbDriverConn, "main")
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize the backup: %v", err)
	}

	isDone, err := backup.Step(-1)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the backup: %v", err)
	}
	if !isDone {
		return nil, fmt.Errorf("Backup is unexpectedly not done")
	}

	err = backup.Finish()
	if err != nil {
		return nil, fmt.Errorf("Failed to finish backup: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("Could not close the tmpFile: %v", err)
	}

	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("Could not read the backup file: %v", err)
	}
	return bytes, nil
}

//TODO: Divide and Conquer
func (d *Datastore) RestoreBackup(content []byte) error {
	//TODO: Testing!!

	tmpFile, err := ioutil.TempFile("", "tempdb")
	if err != nil {
		return fmt.Errorf("Could not create the tempFile for restoring the database: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(content); err != nil {
		return fmt.Errorf("Could not write the content to the tempfile: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("Could not close the tmpFile: %v", err)
	}

	var driverName = fmt.Sprintf("sqlite3_backup_%v", time.Now().UnixNano())
	destFileName := "almue.db"
	// The driver's connection will be needed in order to perform the backup.
	driverConns := []*sqlite3.SQLiteConn{}

	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			driverConns = append(driverConns, conn)
			return nil
		},
	})

	if err := d.Close(); err != nil {
		return fmt.Errorf("Could not close the db for restoring: %v", err)
	}

	if err := os.Remove(destFileName); err != nil {
		return fmt.Errorf("Could not remove the old database: %v", err)
	}

	srcDb, err := sql.Open(driverName, tmpFile.Name())
	if err != nil {
		return fmt.Errorf("Failed to open the source db: %v", err)
	}
	defer srcDb.Close()
	err = srcDb.Ping()
	if err != nil {
		return fmt.Errorf("Failed to connect to the source database: %v", err)
	}

	// Connect to the destination database.
	destDb, err := sql.Open(driverName, destFileName)
	if err != nil {
		return fmt.Errorf("Failed to open the destination database: %v", err)
	}
	defer destDb.Close()
	err = destDb.Ping()
	if err != nil {
		return fmt.Errorf("Failed to connect to the destination database: %v", err)
	}

	if len(driverConns) != 2 {
		return fmt.Errorf("Expected 2 driver connections, but found %v", len(driverConns))
	}
	srcDbDriverConn := driverConns[0]
	if srcDbDriverConn == nil {
		return fmt.Errorf("The source database driver connection is nil")
	}
	destDbDriverConn := driverConns[1]
	if destDbDriverConn == nil {
		return fmt.Errorf("The destination database driver connection is nil")
	}

	backup, err := destDbDriverConn.Backup("main", srcDbDriverConn, "main")
	if err != nil {
		return fmt.Errorf("Failed to initialize the backup: %v", err)
	}

	isDone, err := backup.Step(-1)
	if err != nil {
		return fmt.Errorf("Failed to create the backup: %v", err)
	}
	if !isDone {
		return fmt.Errorf("Backup is unexpectedly not done")
	}

	err = backup.Finish()
	if err != nil {
		return fmt.Errorf("Failed to finish backup: %v", err)
	}

	d.DB, err = sql.Open("sqlite3", "almue.db")
	if err != nil {
		return fmt.Errorf("Could not reopen the restored db")
	}
	return nil
}
