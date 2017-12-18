package sqlu_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/gohxs/sqlu"
	_ "github.com/mattn/go-sqlite3"
)

var (
	now time.Time
)

type User struct {
	ID         string    `sqlu:"id,key"`
	Name       string    `sqlu:"name"`
	CreateTime time.Time `sqlu:"create_date"`
}

func prepareDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "user"
	(
		id integer ,
		name string,
		create_time datetime
	)`)
	if err != nil {
		return nil, err
	}

	now, err = time.Parse("02-01-2006", "18-12-2017")
	if err != nil {
		return nil, err
	}
	// Sample
	_, err = db.Exec(`INSERT INTO "user" VALUES ('1','myname',?)`, now)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestScanNamed(t *testing.T) {
	db, err := prepareDB()
	if err != nil {
		t.Fatal(err)
	}

	res, err := db.Query(`SELECT name FROM "user"`)
	if err != nil {
		t.Fatal(err)
	}
	for res.Next() {
		var user User
		var testUser = User{Name: "myname"}
		err = sqlu.ScanNamed(res, &user)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("User:", user, testUser)
		if !reflect.DeepEqual(user, testUser) {
			t.FailNow()
		}
	}

}

func TestScan(t *testing.T) {
	db, err := prepareDB()
	if err != nil {
		t.Fatal(err)
	}

	res, err := db.Query(`SELECT * FROM "user"`)
	if err != nil {
		t.Fatal(err)
	}
	var user User
	var testUser = User{ID: "1", Name: "myname", CreateTime: now}
	for res.Next() {
		err = sqlu.Scan(res, &user)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("User:", user, testUser)
		if !reflect.DeepEqual(user, testUser) {
			t.FailNow()
		}
	}

}
