package store

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/he4d/simplejack"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "../test.db"

var store *Datastore

func TestMain(m *testing.M) {
	os.Remove(dbPath)

	var err error
	store, err = New(dbPath, simplejack.New(simplejack.TRACE, ioutil.Discard))

	if err != nil {
		log.Fatal(err)
	}
	if store == nil {
		log.Fatal("Store is nil but New didnt return an error")
	}

	code := m.Run()

	clearTable()

	store.Close()

	os.Remove(dbPath)

	os.Exit(code)
}

func TestMigration(t *testing.T) {
	store, err := New(dbPath, simplejack.New(simplejack.TRACE, ioutil.Discard))
	defer store.Close()
	if err != nil {
		t.Error(err)
	}
	if store == nil {
		t.Error("Store is nil but New didnt return an error")
	}
}
