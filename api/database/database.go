package database

import (
	"database/sql"
	"github.com/jenpet/plooral/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	KNoEntityChanged errors.Kind = "NoEntityChanged"
	KNoEntityFound errors.Kind = "NoEntityFound"
)

var dbConnections = map[string]*sql.DB{}

func GetPostgresDatabase(uri string) *sql.DB {
	if con, ok := dbConnections[uri]; ok {
		return con
	}

	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatalf("Could not open connection to postgres database. Error %v", err)
	}
	dbConnections[uri] = db
	return db
}

func GetDefaultPostgresDatabase() *sql.DB {
	return GetPostgresDatabase(parseConfig().PostgresURI)
}

// GenericWrite simplifies an update operation by checking the amount of results and returning an error if no row
// was affected by the change or a general database communication error occurred.
func GenericWrite(db *sql.DB, query string, args ...interface{}) error {
	res, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows <= 0 {
		return errors.E("did not write any entity", KNoEntityChanged, http.StatusBadRequest)
	}
	return err
}

// GenericWriteIDReturn inserts a row into the database and returns the id of the inserted row.
// Condition: The query statement has to have a returning part where it returns the latest (i.e. recently created id)
func GenericWriteIDReturn(db *sql.DB, query string, args...interface{}) (int, error) {
	res := db.QueryRow(query, args...)
	id := -1
	err := res.Scan(&id)
	if err != nil {
		return id, errors.E("did not write any entity", KNoEntityChanged, http.StatusBadRequest, err)
	}
	return id, nil
}

// RowScanner interface allows to use sql.Row and sql.Rows as a scanner having the identical function signature
// sql.Scanner interface does not meet the requirements unfortunately
type RowScanner interface {
	Scan(dest ...interface{}) error
}
