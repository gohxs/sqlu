package sqlu_test

import (
	"reflect"
	"testing"

	"github.com/gohxs/sqlu"
)

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
	var testUser = User{ID: "1", Name: "myname", Alias: "the first", CreateTime: now}
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
