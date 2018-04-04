package sqlu_test

import (
	"database/sql"
	"testing"

	"github.com/gohxs/sqler/sqler"
	"github.com/gohxs/sqlu/x/nonrefl/sqlu"
	"github.com/gohxs/testu/assert"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	_schemaCache *sqlu.Schema
	ID           int    `sqlu:"id,primaryKey"`
	Name         string `sqlu:"name" db:"name"`
	Email        string `sqlu:"email" db:"email"`
}

// Old
func (u User) Table() string {
	return "user"
}

func (m *User) SchemaFields() []interface{} {
	return []interface{}{&m.ID, &m.Name, &m.Email}
}

// Schema
func (m *User) Schema() sqlu.Schema {
	return sqlu.Schema{
		Table: "user",
		Fields: []sqlu.Field{
			{&m.ID, "id", "int"},
			{&m.Name, "name", "string"},
			{&m.Email, "email", "string"},
		},
	}
}

func TestSchema(t *testing.T) {
	a := assert.A(t)
	db, err := sql.Open("sqlite3", ":memory:")
	a.Eq(err, nil, "should not error openeing database")
	_, err = sqlu.Create(db, &User{})
	a.Eq(err, nil, "create user schema")
	_, err = sqlu.Insert(db, &User{
		ID:    1,
		Name:  "Myself",
		Email: "email",
	})
	a.Eq(err, nil, "insert user")

	s := sqler.New()
	s.SetDB(db)
	s.Cmd(`select * from "user" LIMIT 1`)

	u := User{}
	err = sqlu.Get(db, &u, `SELECT * FROM "user" LIMIT 1`)
	a.Eq(err, nil, "should not nil selecting")

	t.Log("User:", u)

}

// Reflection
/*func BenchmarkSQLUOld(b *testing.B) {
	db := prepareDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			u := User{}
			row, err := db.Query(`SELECT name FROM "user" LIMIT 1`)
			if err != nil {
				b.Fatal(err)
			}
			defer row.Close()
			if !row.Next() {
				b.Fatal("no rows")
			}
			err = sqluold.Scan(row, &u)
			if err != nil {
				b.Fatal(err)
			}
			if u.ID != 0 && u.Email != "" {
				b.Fatal("Error should be empty")
			}
		}()
	}

}*/
var tQry = `SELECT name FROM "user"`

func BenchmarkSQLX(b *testing.B) {
	db := sqlx.NewDb(prepareDB(b), "sqlite3")
	//db := prepareDB(b)

	b.ResetTimer()
	rows, err := db.Queryx(tQry)
	u := User{}
	for i := 0; i < b.N; i++ {
		if !rows.Next() {
			rows.Close()
			rows, err = db.Queryx(tQry)
			if err != nil {
				b.Fatal(err)
			}
			continue
		}
		err = rows.StructScan(&u)

		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkSQLU(b *testing.B) {
	db := prepareDB(b)
	b.ResetTimer()
	rows, err := db.Query(tQry)
	u := User{}
	rowScan := sqlu.NewRowScan(rows)
	for i := 0; i < b.N; i++ {
		if !rows.Next() {
			rows.Close()
			rows, err = db.Query(tQry)
			if err != nil {
				b.Fatal(err)
			}
			rowScan = sqlu.NewRowScan(rows)
			continue
		}
		err = rowScan.Scan(&u)
		if err != nil {
			b.Fatal(err)
		}

	}

}

func prepareDB(b *testing.B) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	_, err = sqlu.Create(db, &User{})
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 100000; i++ { // Insert 100
		_, err = sqlu.Insert(db, &User{ID: 1, Name: "test", Email: "test@test.t"})
		if err != nil {
			b.Fatal(err)
		}
	}

	return db
}
