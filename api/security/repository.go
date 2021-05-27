package security

import (
	"database/sql"
	"github.com/jenpet/plooral/database"
)

func newRepository(uri string) *repository {
	return &repository{db: database.GetPostgresDatabase(uri)}
}

func newDefaultRepository() *repository {
	return &repository{db: database.GetDefaultPostgresDatabase()}
}

type repository struct {
	db *sql.DB
}

func (r *repository) persistPassword(password string) (*CredentialSet, error) {
	q := `
		INSERT INTO
			security_credentials( password )
		VALUES(
			crypt($1, gen_salt('bf'))
		)
		RETURNING
			id;
	`
	id, err := database.GenericWriteIDReturn(r.db, q, password)
	if err != nil {
		return nil, err
	}
	return &CredentialSet{ ID: id, Password: password }, nil
}

func (r *repository) verifyPassword(id int, password string) (bool, error) {
	q := `
		SELECT
			COUNT(1)
		FROM
			security_credentials
		WHERE
			id = $1
		AND
			active = true
		AND
			password = crypt($2, password);
		`
	row := r.db.QueryRow(q, id, password)
	amount := -1
	err := row.Scan(&amount)
	return amount > 0, err
}