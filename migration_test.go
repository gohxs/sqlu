package sqlu_test

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/gohxs/prettylog"
	"github.com/gohxs/sqler/sqler"
	"github.com/gohxs/sqlu"
	_ "github.com/mattn/go-sqlite3"
)

var (
	migrators = []sqlu.M{
		initial,
	}
)

var (
	initial = sqlu.M{
		Name: "initial",
		Up:   `CREATE TABLE "my" (id serial primary key, name string)`,
		Down: `DROP TABLE "my"`,
	}
)

func init() {
	prettylog.Global()
}

func TestMain(t *testing.T) {
	db, err := sql.Open("sqlite3", "tmp.sqlite3")
	mig, err := sqlu.NewMigrator(db, "_migrations")
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Run(migrators)
	if err != nil {
		t.Fatal(err)
	}
	s := sqler.New()
	s.SetDB(db)
	s.Cmd(`select * from "_migrator"`)

}

func TestMigrationError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	mig, err := sqlu.NewMigrator(db, "_migrations")
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Run([]sqlu.M{
		{
			Name: "initial",
			Up:   `CREATE TABLE mymigrator (id serial)`,
			Down: `DROP TABLE mymigrator`,
		},
		{
			Name: "testerror",
			Up:   func(tx *sql.Tx) error { return errors.New("empty") },
		},
	})
	if err != nil {
		log.Println("Migration err:", err)
	}

	s := sqler.New()
	s.SetDB(db)
	s.Cmd(`select * from "_migrations"`)

}
