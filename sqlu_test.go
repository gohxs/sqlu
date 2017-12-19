package sqlu_test

import (
	"database/sql"
	"log"
	"os"
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
	Alias      string    `sqlu:"nick,omitempty"`
	Age        int       `sqlu:"age"`
	CreateTime time.Time `sqlu:"create_date,createTimeStamp"`
	UpdateTime time.Time `sqlu:"update_date,updateTimeStamp"`
	NonSQLU    int
}

func (u *User) Table() string { return "user" }

func init() {
	sqlu.Log = log.New(os.Stderr, "", log.LstdFlags)
}

func prepareDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = sqlu.Create(db, &User{})
	if err != nil {
		t.Fatal(err)
	}

	user := User{ID: "1", Name: "myname", Alias: "the first", Age: 1}
	_, err = sqlu.Insert(db, &user)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
