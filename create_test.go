package sqlu_test

import (
	"testing"
	"time"

	"github.com/gohxs/sqlu"
)

type Brand struct {
	ID      int       `sqlu:"id"`
	Name    string    `sqlu:"name"`
	Created time.Time `sqlu:"create_date,,createTimeStamp"`
}

func (b *Brand) Table() string { return "brand" }

func TestCreate(t *testing.T) {
	db := prepareDB(t)

	_, err := sqlu.Create(db, &Brand{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = sqlu.Insert(db, &Brand{ID: 1, Name: "Google"})
	if err != nil {
		t.Fatal(err)
	}

	res, err := db.Query("SELECT * FROM BRAND")
	if err != nil {
		t.Fatal(err)
	}

	for res.Next() {
		var b Brand
		sqlu.Scan(res, &b)

		t.Logf("R: %v", b)

	}

}
