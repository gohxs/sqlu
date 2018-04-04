package sqlu

import (
	"database/sql"
)

// Helper
func Transact(db DBer, fn func(tx *sql.Tx) error) error {
	var err error
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = fn(tx)
	return err
}
