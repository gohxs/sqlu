package sqlu_test

import (
	"testing"

	"github.com/gohxs/sqlu"
)

func TestDelete(t *testing.T) {
	db := prepareDB(t)
	var err error
	_, err = sqlu.Insert(db, &User{ID: "test", Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	// Select test
	{
		res, err := db.Query(`SELECT * FROM "user"`)
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for res.Next() {
			count++
		}
		if count != 2 {
			t.Fatal("Failed test insertion")
		}
	}

	res, err := sqlu.Delete(db, &User{ID: "test"})
	if err != nil {
		t.Fatal(err)
	}

	nrows, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if nrows != 1 {
		t.Fatal("Affected rows should be 1")
	}

	// Query again
	{
		res, err := db.Query(`SELECT * FROM "user"`)
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for res.Next() {
			count++
		}
		if count != 1 {
			t.Fatal("Failed test deletion")
		}
	}

}
