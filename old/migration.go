package sqlu

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// Migration states
const (
	StateErr State = -1 + iota
	StateNil
	StateStarted
	StateFailed
	StateFinished
)

// State migration state
type State int

func (s State) String() string {
	migStr := map[State]string{
		StateErr:      "Unknown",
		StateNil:      "No migration",
		StateStarted:  "started",
		StateFailed:   "failed",
		StateFinished: "finished",
	}
	return migStr[s]
}

// MigratorFunc
// type MigratorFunc func(tx *sql.Tx) error

// M migration type
type M struct {
	Name string
	Up   interface{}
	Down interface{}
}

// Manager migration manager
type Manager struct {
	db      *sql.DB
	tblName string
}

// New creates a migration manager
func New(db *sql.DB, tblName string) (*Manager, error) {
	m := &Manager{db, tblName}
	if err := m.Init(); err != nil {
		return nil, err
	}
	return m, nil
}

// Init initializes migration datatables
func (m *Manager) Init() error {
	// Depends on database
	log.Println("Initialize migration table")
	// Create database schema for migrations
	_, err := m.db.Exec(`
	CREATE TABLE IF NOT EXISTS "` + m.tblName + `"(
		id integer primary key autoincrement, 
		name text,
		state integer,
		created_at datetime
	)
	`)
	return err
}

// Run will run a migration list
func (m *Manager) Run(migs []M) error {
	log.Printf("Running '%d' migrations", len(migs))

	for _, mig := range migs {
		state, err := m.GetMigrationState(mig.Name)
		if err != nil {
			return err
		}
		log.Printf("Current state for migration '%s' --- '%s'", mig.Name, state.String())
		if state == StateFinished {
			log.Println("Migration exists, continuing")
			continue
		}

		err = m.StartMigration(mig)
		if err != nil {
			return err
		}
	}
	return nil
}

// StartMigration migrates UP
func (m *Manager) StartMigration(mig M) error {
	err := m.SaveMigration(mig.Name, StateStarted)
	if err != nil {
		return err
	}

	log.Printf("Running migration: '%s'", mig.Name)

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	err = func() error {
		switch v := mig.Up.(type) {
		case string:
			_, err := tx.Exec(v)
			return err
		case func(tx *sql.Tx) error:
			log.Println("Migration func")
			return v(tx)
		}
		return errors.New("Unsupported migration type")
	}()

	if err != nil {
		log.Println("FAIL:", err)
		tx.Rollback()
		m.SaveMigration(mig.Name, StateFailed)
		return err
	}

	tx.Commit()
	return m.SaveMigration(mig.Name, StateFinished)
}

// GetMigrationState fetches a migration state
func (m *Manager) GetMigrationState(name string) (State, error) {
	res := m.db.QueryRow(`
		SELECT state FROM "`+m.tblName+`" 
		WHERE name=$1
		`, name)

	state := StateErr
	err := res.Scan(&state)
	if err != nil && err != sql.ErrNoRows {
		return StateErr, err
	}
	if err == sql.ErrNoRows {
		return StateNil, nil
	}
	return state, nil
}

// SaveMigration update or insert a migration
func (m *Manager) SaveMigration(name string, state State) error {
	log.Printf("MIGRATION '%s' state -> '%s'", name, state)
	dbState, err := m.GetMigrationState(name)
	if err != nil {
		log.Println("Err:", err)
		return err
	}
	if dbState == StateNil { // We Insert
		_, err := m.db.Exec(`
			INSERT INTO "`+m.tblName+`" 
			(name,state,created_at)
			values($1,$2, $3)`, name, state, time.Now().UTC())
		if err != nil {
			return err
		}
	} else {
		_, err := m.db.Exec(`
		UPDATE "`+m.tblName+`"
		SET state = $1
		WHERE name = $2`, state, name)
		if err != nil {
			return err
		}
	}
	return nil
}
