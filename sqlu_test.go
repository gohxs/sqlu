package sqlu_test

import (
	"database/sql"
	"testing"

	"github.com/gohxs/sqler/sqler"
	"github.com/gohxs/sqlu"
	"github.com/gohxs/testu/assert"
	"github.com/jmoiron/sqlx"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

type Common struct {
	Test string
}

type User struct {
	Common
	ID    int    `sqlu:"id,primaryKey"`
	Name  string `sqlu:"name" db:"name"`
	Email string `sqlu:"email" db:"email"`

	Address string
	Field1  string
	Field2  string
	Another string
	Test    struct{ Name string }
}

// Old

// Schema
func (m *User) Schema(s *sqlu.S) {
	s.BuildSchema("user", m, sqlu.ReflectSchema, m.Fields)
}
func (m *User) Fields(values []interface{}, indexes ...int) {
	sqlu.ReflectFieldsx(m, values, indexes...)
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
		rows, err := sqlu.NewRowScanner(db.Query(tQry))
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := User{}
			//fields := u.Schema().Fields()
			if !rows.Next() {
				rows.Close()
				rows, err = sqlu.NewRowScanner(db.Query(tQry))
				if err != nil {
					b.Fatal(err)
				}
				continue
			}
			//err = rows.Scan(fields...)
			err = rows.Scan(&u)
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
		rows, err := sqlu.NewRowScanner(db.Query(tQry))
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			u := User{}
			//fields := u.Schema().Fields()
			if !rows.Next() {
				rows.Close()
				rows, err = sqlu.NewRowScanner(db.Query(tQry))
				if err != nil {
					b.Fatal(err)
				}
				continue
			}
			//err = sqlu.Scan(rows, &u) //.Scan(fields...)
			err = rows.Scan(&u)
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
	for i := 0; i < 10000; i++ { // Insert 100
		u := &User{ID: 1, Name: "test", Email: "test@test.t"}
		sqlu.Insert(db, u)
	}

	return db
}
