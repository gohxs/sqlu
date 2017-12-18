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
	_, err = db.Exec(`INSERT INTO "user" VALUES ('1','myname',?)`, now)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestScanNamed(t *testing.T) {
	db := prepareDB(t)

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
	db := prepareDB(t)

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

func TestTableInsert(t *testing.T) {
	db := prepareDB(t)

	res, err := sqlu.TableInsert(db, "user", &User{ID: "2", Name: "name2", CreateTime: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	if v, err := res.RowsAffected(); v != 1 && err != nil {
		t.Fatal("Rows affected should be 1")
	}
}

func TestInsert(t *testing.T) {
	db := prepareDB(t)

	res, err := sqlu.Insert(db, &User{ID: "2", Name: "name2", CreateTime: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	if v, err := res.RowsAffected(); v != 1 && err != nil {
		t.Fatal("Rows affected should be 1")
	}
}
