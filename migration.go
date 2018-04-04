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

// M migration operations
type M struct {
	Name string
	Up   interface{}
	Down interface{}
}

// Migrator migration manager
type Migrator struct {
	db      *sql.DB
	tblName string
}

type MigrationModel struct {
	Table     string
	ID        int
	Name      string
	State     int
	CreatedAt time.Time
}

func initMigrationSchema(s *Schema) {
	s.
		Field("id", "integer primary key autoincrement").
		Field("name", "text").
		Field("state", "int").
		Field("created_at", "datetime")

}
func (m *MigrationModel) Schema() *Schema {
	return BuildSchema(
		m.Table,
		initMigrationSchema,
	)
}
func (m *MigrationModel) Fields() []interface{} {
	return Fields(&m.ID, &m.Name, &m.State, &m.CreatedAt)
}

// New creates a migration manager
func NewMigrator(db *sql.DB, tblName string) (*Migrator, error) {
	m := &Migrator{db, tblName}
	if err := m.Init(); err != nil {
		return nil, err
	}
	return m, nil
}

// Init initializes migration datatables
func (m *Migrator) Init() error {
	_, err := Create(m.db, &MigrationModel{Table: m.tblName})
	return err
}

// Run will run a migration list
func (m *Migrator) Run(migs []M) error {
	log.Printf("Running '%d' migrations", len(migs))

	for _, mig := range migs {
		state, err := m.GetMigrationState(mig.Name)
		if err != nil {
			return err
		}
		if state == StateFinished {
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
func (m *Migrator) StartMigration(mig M) error {
	err := m.SaveMigration(mig.Name, StateStarted)
	if err != nil {
		return err
	}
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
		tx.Rollback()
		m.SaveMigration(mig.Name, StateFailed)
		return err
	}

	tx.Commit()
	return m.SaveMigration(mig.Name, StateFinished)
}

// GetMigrationState fetches a migration state
func (m *Migrator) GetMigrationState(name string) (State, error) {
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
func (m *Migrator) SaveMigration(name string, state State) error {
	dbState, err := m.GetMigrationState(name)
	if err != nil {
		log.Println("Err:", err)
		return err
	}
	if dbState == StateNil { // We Insert
		mm := MigrationModel{
			Table:     m.tblName,
			Name:      name,
			State:     int(state),
			CreatedAt: time.Now().UTC(),
		}
		_, err := Insert(m.db, &mm)
		if err != nil {
			return err
		}
	} else {
		Update(m.db, &MigrationModel{Table: m.tblName, State: int(state)}, []string{"state"}, "name = ?", name)
		if err != nil {
			return err
		}
	}
	return nil
}
