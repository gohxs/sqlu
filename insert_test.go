package sqlu_test

import (
	"testing"
	"time"

	"github.com/gohxs/sqlu"
)

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
