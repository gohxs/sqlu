package sqlu_test

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	now time.Time
)

type User struct {
	ID         string    `sqlu:"id,key"`
	Name       string    `sqlu:"name"`
	Alias      string    `sqlu:"nick"`
	CreateTime time.Time `sqlu:"create_date,,createTimeStamp"`
}

func (u *User) Table() string { return "user" }

func prepareDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "user"
	(
		id integer ,
		name string,
		nick string,
		create_date datetime
	)`)
	if err != nil {
		t.Fatal(err)
	}

	now, err = time.Parse("02-01-2006", "18-12-2017")
	if err != nil {
		t.Fatal(err)
	}
	// Sample
	_, err = db.Exec(`INSERT INTO "user" VALUES ('1','myname','the first',?)`, now)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
