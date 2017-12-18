package sqlu_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/gohxs/sqlu"
	_ "github.com/mattn/go-sqlite3"
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

	// Sample
	_, err = db.Exec(`INSERT INTO "user" VALUES ('1','myname',?)`, time.Now())
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
		err = sqlu.ScanNamed(res, &user)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("User:", user)
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
	for res.Next() {
		err = sqlu.Scan(res, &user)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("User:", user)
	}

}
