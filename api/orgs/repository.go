package orgs

import (
	"database/sql"
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/security"
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
			protected,
		    user_credentials_id,
			owner_credentials_id
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

func (r *repository) organizationExists(slug string) (bool, error) {
	q := `
		SELECT 
			COUNT(1)
		FROM
			organizations
		WHERE
			slug = $1
	`
	res := r.db.QueryRow(q, slug)
	amount := 0
	err := res.Scan(&amount)
	return amount > 0, err
}

func (r *repository) organizationBySlug(slug string) (*Organization, error) {
	q := `
		SELECT
			id,
			slug,
			title,
			description,
			hidden,
			protected,
		    user_credentials_id,
			owner_credentials_id
		FROM
			organizations
		WHERE
			slug = $1
	`
	return r.processOrganizationRow(r.db.QueryRow(q, slug))
}

func (r *repository) createOrganization(o Organization) (*Organization, error) {
	q := `
		INSERT INTO
			organizations (slug, title, description, hidden, protected, user_credentials_id, owner_credentials_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var userCredentialID *int
	if o.UserSecurity.HasValidID() {
		userCredentialID = &o.UserSecurity.ID
	}
	var ownerCredentialID *int
	if o.OwnerSecurity.HasValidID() {
		ownerCredentialID = &o.OwnerSecurity.ID
	}
	id, err := database.GenericWriteIDReturn(
		r.db,
		q,
		o.Slug,
		o.Name,
		o.Description,
		o.Hidden,
		o.Protected,
		userCredentialID,
		ownerCredentialID)
	if err != nil {
		return nil, err
	}
	o.ID = id
	return &o, nil
}

func (r *repository) updateOrganization(o Organization) (*Organization, error) {
	q := `
		UPDATE
			organizations
		SET
			title = $3,
			description = $4,
			hidden = $5,
			protected = $6,
			modified = (NOW() AT TIME ZONE 'utc')
		WHERE
			id = $1
		AND
			slug = $2
		RETURNING id
	`
	id, err := database.GenericWriteIDReturn(r.db, q, o.ID, o.Slug, o.Name, o.Description, o.Hidden, o.Protected)
	if err != nil {
		return nil, err
	}
	o.ID = id
	return &o, nil
}

func (r *repository) processOrganizationRow(sc database.RowScanner) (*Organization, error) {
	var dbo dbOrganziation
	err := sc.Scan(
		&dbo.ID,
		&dbo.Slug,
		&dbo.Name,
		&dbo.Description,
		&dbo.Hidden,
		&dbo.Protected,
		&dbo.userCredentialsID,
		&dbo.ownerCredentialsID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Ef("no organization found", database.KNoEntityFound, http.StatusNotFound)
		}
		return nil, err
	}
	o := dbo.toOrganization()
	return &o, nil
}

type dbOrganziation struct {
	Organization
	userCredentialsID, ownerCredentialsID sql.NullInt32
}

func (dbo dbOrganziation) toOrganization() Organization {
	o := dbo.Organization
	if dbo.userCredentialsID.Valid {
		userSecurity := security.CredentialSet{ ID: int(dbo.userCredentialsID.Int32) }
		o.UserSecurity = &userSecurity
	}

	if dbo.ownerCredentialsID.Valid {
		ownerSecurity := security.CredentialSet{ ID: int(dbo.ownerCredentialsID.Int32) }
		o.OwnerSecurity = &ownerSecurity
	}
	return o
}
