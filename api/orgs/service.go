package orgs

import (
	"context"
	"fmt"
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"github.com/jenpet/plooral/security"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const KOrganizationStateCorrupt errors.Kind = "OrganizationStateCorrupt"

type Service struct {
	repo *repository
	cs   credentialService
}

func newDefaultService(ps credentialService) *Service {
	return newService(newDefaultRepository(), ps)
}

func newService(r *repository, ps credentialService) *Service {
	return &Service{repo: r, cs: ps}
}

func (s *Service) CreateOrganization(ctx context.Context, o partialOrganization) (*Organization, error) {
	// hidden and unprotected organizations cannot be created and are considered invalid
	if o.isHidden() && !o.isProtected() {
		return nil, errors.Ef("hidden organizations also have to be protected", http.StatusBadRequest, rest.KUserInputInvalid)
	}

	if o.isProtected() && o.PartialCredentialSet == nil {
		return nil, errors.Ef("protected organizations require a password", http.StatusBadRequest, rest.KUserInputInvalid)
	}

	// check if organization exists
	exists, err := s.repo.organizationExists(*o.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.Ef("organization '%s' does already exist.", *o.Slug, http.StatusBadRequest, rest.KUserInputInvalid)
	}

	newOrg := o.toOrganization()

	// in case the organization is protected generate the user and owner credentials and set them accordingly
	if o.isProtected() {
		// persist the credential set given by the user
		userSet, err := s.cs.PersistCredentials(*o.PartialCredentialSet)
		if err != nil {
			return nil, errors.Ef("failed to persist credentials. Error: %+v", err)
		}

		// create and persist credential set for the owner
		ownerSet, err := s.cs.GenerateAndPersistPassword()
		if err != nil {
			return nil, errors.Ef("failed to persist credentials. Error: %+v", err)
		}
		newOrg.UserSecurity = userSet
		newOrg.OwnerSecurity = ownerSet
	}

	created, err := s.repo.createOrganization(newOrg)
	return created, err
}

func (s *Service) UpdateOrganization(ctx context.Context, o partialOrganization) (*Organization, error) {
	org, err := s.OrganizationBySlug(ctx, *o.Slug)
	if err != nil {
		return nil, err
	}
	org.mergeWithPartial(o)
	return s.repo.updateOrganization(*org)
}

func (s *Service) OrganizationBySlug(ctx context.Context, slug string) (*Organization, error) {
	org, err := s.repo.organizationBySlug(slug)
	if err != nil {
		return nil, err
	}
	// if the organization is not protected simply return it
	if !org.Protected {
		org.clearCredentials()
		return org, err
	}
	creds := extractCredentials(ctx)
	if err = s.validateOrganizationAccess(*org, creds); err != nil {
		return nil, err
	}
	org.clearCredentials()
	return org, nil
}

func (s *Service) validateOrganizationAccess(org Organization, creds string) error {
	// check if th org at least a valid owner security
	if !org.UserSecurity.HasValidID() || !org.OwnerSecurity.HasValidID() {
		return errors.Ef("organization '%s' is in a corrupted state, protected but no credentials present.", KOrganizationStateCorrupt)
	}
	// in case user security and owner security do not match an error will be returned
	// TODO: Maybe remove calls to the credential service to reduce database interactions and perform joins with org table
	if !s.validateCredentials(org.UserSecurity.ID, creds) && !s.validateCredentials(org.OwnerSecurity.ID, creds) {
		// hidden orgs with wrong credentials should result in a 404 NOT FOUND HTTP response code to avoid reconstructing hidden org structures
		if org.Hidden {
			return errors.E("no organization found", database.KNoEntityFound, http.StatusNotFound)
		}
		return errors.E("invalid credentials", security.KCredentialInputInvalid, http.StatusForbidden)
	}
	return nil
}

func (s *Service) validateCredentials(credentialSetID int, creds string) bool {
	ok, err := s.cs.VerifyCredentials(credentialSetID, creds)

	if err != nil {
		log.Warnf("failed verifying credentials for credential set ID '%d'", credentialSetID)
		return false
	}
	return ok
}

func (s *Service) AllOrganizations(includeHidden bool) ([]Organization, error) {
	orgs, err := s.repo.allOrganizations(includeHidden)
	if err != nil {
		return []Organization{}, nil
	}
	for i := range orgs {
		orgs[i].clearCredentials()
	}
	return orgs, err
}

func extractCredentials(ctx context.Context) string {
	creds := ctx.Value(rest.CredentialContextKey)
	if creds == nil {
		return ""
	}
	return fmt.Sprintf("%s", creds)
}

type credentialService interface {
	GenerateAndPersistPassword() (*security.CredentialSet, error)
	VerifyCredentials(id int, password string) (bool, error)
	PersistCredentials(set security.PartialCredentialSet) (*security.CredentialSet, error)
}