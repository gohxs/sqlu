package sqlu_test

import (
	"database/sql"
	"testing"

	"github.com/gohxs/sqlu"
	"github.com/gohxs/testu/assert"
)

func TestFind(t *testing.T) {
	a := assert.A(t)
	db, err := sql.Open("sqlite3", ":memory:")
	a.Eq(err, nil, "error should be nil")

	_, err = sqlu.Create(db, &User{})
	a.Eq(err, nil, "should create a table")

	_, err = sqlu.Find(db,
		sqlu.Sample(&User{Name: "hello"}, "name", "email"),
		sqlu.Sample(&User{Name: "test"}, "name"),
	)
	a.Eq(err, nil, "Should find entries")

}
