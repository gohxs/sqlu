package sqlu_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/gohxs/sqler/sqler"
	"github.com/gohxs/sqlu"
	"github.com/gohxs/testu/assert"
	"github.com/jmoiron/sqlx"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `sqlu:"id,primaryKey"`
	Name  string `sqlu:"name" db:"name"`
	Email string `sqlu:"email" db:"email"`
}

// Old

// Schema
func (m *User) Schema() *sqlu.Schema {
	return sqlu.BuildSchema(
		"user",
		func(s *sqlu.Schema) {
			s.
				Field("id", "int").
				Field("name", "text").
				Field("email", "text")
		},
	)
}
func (m *User) Fields() []interface{} {
	return sqlu.Fields(&m.ID, &m.Name, &m.Email)
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

	//u := User{}
	//err = sqlu.Get(db, &u, `SELECT * FROM "user" LIMIT 1`)
	//a.Eq(err, nil, "should not nil selecting")

	//t.Log("User:", u)

}

func BenchmarkAsterisk(b *testing.B) {
	var tQry = `SELECT * FROM "user"`

	b.Run("SQLX", func(b *testing.B) {
		dbx := sqlx.NewDb(prepareDB(b), "sqlite3")
		rows, err := dbx.Queryx(tQry)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := User{}
			if !rows.Next() {
				rows.Close()
				rows, err = dbx.Queryx(tQry)
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
	})

	b.Run("SQLU", func(b *testing.B) {
		db := prepareDB(b)
		rows, err := db.Query(tQry)
		if err != nil {
			b.Fatal(err)
		}
		rowScan := sqlu.NewRowScanner(rows)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := User{}
			//fields := u.Schema().Fields()
			if !rows.Next() {
				rows.Close()
				rows, err = db.Query(tQry)
				if err != nil {
					b.Fatal(err)
				}
				rowScan = sqlu.NewRowScanner(rows)
				continue
			}
			//err = rows.Scan(fields...)
			err = rowScan.Scan(&u)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
func BenchmarkField(b *testing.B) {
	var tQry = `SELECT name,email FROM "user"`

	b.Run("SQLX", func(b *testing.B) {
		dbx := sqlx.NewDb(prepareDB(b), "sqlite3")
		rows, err := dbx.Queryx(tQry)
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := User{}
			if !rows.Next() {
				rows.Close()
				rows, err = dbx.Queryx(tQry)
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
	})

	b.Run("SQLU", func(b *testing.B) {
		db := prepareDB(b)
		rows, err := db.Query(tQry)
		if err != nil {
			b.Fatal(err)
		}
		rowScan := sqlu.NewRowScanner(rows)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := User{}
			//fields := u.Schema().Fields()
			if !rows.Next() {
				rows.Close()
				rows, err = db.Query(tQry)
				if err != nil {
					b.Fatal(err)
				}
				rowScan = sqlu.NewRowScanner(rows)
				continue
			}
			//err = rows.Scan(fields...)
			err = rowScan.Scan(&u)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func prepareDB(b *testing.B) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	db.Exec(`DROP TABLE "user"`)
	_, err = sqlu.Create(db, &User{})
	if err != nil {
		b.Fatal(err)
	}
	values := []string{}
	for i := 0; i < 10000; i++ { // Insert 100
		u := &User{ID: 1, Name: "test", Email: "test@test.t"}
		values = append(values, fmt.Sprintf("(%d,'%s','%s')", u.ID, u.Name, u.Email))
	}
	fullQry := `INSERT INTO "user" (id,name,email) values` + strings.Join(values, ",")

	r, err := db.Exec(fullQry)
	if err != nil {
		b.Fatal(err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected != 10000 {
		b.Fatal(fmt.Errorf("Not enough rows"))
	}

	return db
}
