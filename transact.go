package sqlu

// Helper
func Transact(db Queryer, fn func(tx TxQueryer) error) error {
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
