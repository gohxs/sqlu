package sqlu

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

//Errors
var (
	ErrNotPointer = errors.New("param is not a pointer")
)

const (
	tagField     = 0
	tagOmitEmpty = 1
	tagType      = 2
)

type tagOptions struct {
	fieldName       string
	PrimaryKey      bool
	OmitEmpty       bool
	CreateTimeStamp bool
	UpdateTimeStamp bool
	Unique          bool
	NotNull         bool
}

func parseTag(tagStr string) *tagOptions {
	ret := &tagOptions{}
	tags := strings.Split(tagStr, ",")

	if len(tags) == 0 {
		return ret
	}
	ret.fieldName = tags[0]
	for _, t := range tags[1:] {
		switch strings.ToLower(t) {
		case "primarykey":
			ret.PrimaryKey = true
		case "omitempty":
			ret.OmitEmpty = true
		case "createtimestamp":
			ret.CreateTimeStamp = true
		case "updatetimestamp":
			ret.UpdateTimeStamp = true
		case "unique":
			ret.Unique = true
		case "notnull":
			ret.NotNull = true
		}
	}
	return ret
}

// SQLer interface for "sql" package
// Can be sql.DB or sql.Tx
type SQLer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Tabler Interface to return the table name of a struct
type Tabler interface {
	Table() string
}
