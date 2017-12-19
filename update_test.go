package sqlu_test

import (
	"context"
	"testing"

	"github.com/gohxs/sqlu"
)

func TestUpdate(t *testing.T) {
	db := prepareDB(t)

	var err error
	_, err = sqlu.Insert(db, &User{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = sqlu.Insert(db, &User{ID: "3", Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	var updUser = User{ID: "3", Name: "Different", Alias: "Alias"}
	sqlu.TableUpdateContext(context.Background(), db, "user", &updUser)

	res, err := db.Query("SELECT * FROM USER")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Close()
	for res.Next() {
		var u User
		sqlu.Scan(res, &u)
		t.Log("User: ", u)
	}
}

func TestOmitEmpty(t *testing.T) {
	db := prepareDB(t)
	var updUser = User{ID: "1"}

	sqlu.Update(db, &updUser)

	res, err := db.Query(`SELECT * FROM "user" WHERE id = 1`)
	if err != nil {
		t.Fatal(err)
	}
	for res.Next() {
		var user User
		sqlu.Scan(res, &user)
		if user.Name != "" {
			t.Fatal("Name should be empty")
		}
		if user.Alias == "" {
			t.Fatal("Alias should not be empty [omitempty]")
		}
	}

}
