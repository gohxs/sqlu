package sqlu_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/gohxs/sqlu"
	_ "github.com/mattn/go-sqlite3"
)

var (
	now time.Time
)

type User struct {
	ID         string    `sqlu:"id,primaryKey"`
	Name       string    `sqlu:"name"`
	Alias      string    `sqlu:"nick"`
	CreateTime time.Time `sqlu:"create_date,createTimeStamp"`
	UpdateTime time.Time `sqlu:"update_date,updateTimeStamp"`
}

func (u *User) Table() string { return "user" }

func prepareDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	sqlu.Create(db, &User{})
	if err != nil {
		t.Fatal(err)
	}

	user := User{ID: "1", Name: "myname", Alias: "the first"}
	sqlu.Insert(db, &user)

	return db
}
