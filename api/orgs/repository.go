package orgs

import (
	"database/sql"
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"net/http"
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

func (r *repository) allOrganizations(includeHidden bool) ([]Organization, error) {
	var orgs = make([]Organization, 0)
	q := `
		SELECT
			id,
			slug,
			title,
			description,
			hidden,
			protected
		FROM
			organizations
	`
	// filter out hidden orgs from the resultset
	if !includeHidden {
		q += `
		WHERE
			hidden=false
	`
	}
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		o, err := r.processOrganizationRow(rows)
		if err != nil {
			return nil, err
		}

		if includeHidden || !o.Hidden {
			orgs = append(orgs, *o)
		}
	}
	return orgs, nil
}

func (r *repository) organizationBySlug(slug string) (*Organization, error) {
	q := `
		SELECT
			id,
			slug,
			title,
			description,
			hidden,
			protected
		FROM
			organizations
		WHERE
			slug = $1
	`
	return r.processOrganizationRow(r.db.QueryRow(q, slug))
}

func (r *repository) upsertOrganization(o Organization) (*Organization, error) {
	q := `
		INSERT INTO
			organizations (slug, title, description, hidden, protected)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT(slug)
		DO
			UPDATE SET
				title = $2,
				description = $3,
				hidden = $4,
				protected = $5,
				modified = (NOW() AT TIME ZONE 'utc')
		RETURNING id
	`
	id, err := database.GenericWriteIDReturn(r.db, q, o.Slug, o.Name, o.Description, o.Hidden, o.Protected)
	if err != nil {
		return nil, err
	}
	o.ID = id
	return &o, nil
}

func (r *repository) processOrganizationRow(sc database.RowScanner) (*Organization, error) {
	var o Organization
	err := sc.Scan(
		&o.ID,
		&o.Slug,
		&o.Name,
		&o.Description,
		&o.Hidden,
		&o.Protected)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Ef("no organization found", database.KNoEntityFound, http.StatusNotFound)
		}
		return nil, err
	}
	return &o, nil
}
