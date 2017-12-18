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
func TestCreateTimestamp(t *testing.T) {
	db := prepareDB(t)

	now := time.Now()

	res, err := sqlu.Insert(db, &User{ID: "2", Name: "name2"})
	if err != nil {
		t.Fatal(err)
	}
	if v, err := res.RowsAffected(); v != 1 && err != nil {
		t.Fatal("Rows affected should be 1")
	}

	{
		res, err := db.Query(`SELECT * FROM "user" WHERE id = '2'`)
		if err != nil {
			t.Fatal(err)
		}

		if res.Next() {
			var user User
			sqlu.Scan(res, &user)

			t.Logf("R: %#v %#v", user, now)
			if user.CreateTime.Day() != now.Day() {
				t.Fatal("Day is not equal")
			}
		}
	}

}
