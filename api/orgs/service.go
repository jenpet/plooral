package orgs

import (
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"github.com/jenpet/plooral/security"
	"net/http"
)

type Service struct {
	repo *repository
	ps passwordService
}

func newDefaultService(ps passwordService) *Service {
	return newService(newDefaultRepository(), ps)
}

func newService(r *repository, ps passwordService) *Service {
	return &Service{repo: r, ps: ps}
}

func (s *Service) CreateOrganization(o partialOrganization) (*Organization, error) {
	// hidden and unprotected organizations cannot be created and are considered invalid
	if o.isHidden() && !o.isProtected() {
		return nil, errors.Ef("hidden organizations also have to be protected", http.StatusBadRequest, rest.KUserInputInvalid)
	}

	if o.isProtected() && o.PartialCredentialSet == nil {
		return nil, errors.Ef("protected organizations require a password", http.StatusBadRequest, rest.KUserInputInvalid)
	}

	// check if organization exists
	org, err := s.OrganizationBySlug(*o.Slug)
	if err != nil && errors.ErrKind(err) != database.KNoEntityFound {
		return nil, err
	}
	if org != nil {
		return nil, errors.E("organization '%s' does already exist.", http.StatusBadRequest, rest.KUserInputInvalid)
	}

	// persist the credential set given by the user
	userSet, err := s.ps.PersistCredentials(*o.PartialCredentialSet)
	if err != nil {
		return nil, errors.Ef("failed to persist credentials. Error: %+v", err)
	}

	// create and persist credential set for the owner
	ownerSet, err := s.ps.GenerateAndPersistPassword()
	if err != nil {
		return nil, errors.Ef("failed to persist credentials. Error: %+v", err)
	}

	// create organization
	newOrg := o.toOrganization()
	newOrg.UserSecurity = userSet
	newOrg.OwnerSecurity = ownerSet

	created, err := s.repo.upsertOrganization(newOrg)
	return created, err
}

func (s *Service) UpdateOrganization(o partialOrganization) (*Organization, error) {
	org, err := s.OrganizationBySlug(*o.Slug)
	if err != nil {
		return nil, err
	}
	org.mergeWithPartial(o)
	return s.repo.upsertOrganization(*org)
}

func (s *Service) OrganizationBySlug(slug string) (*Organization, error) {
	return s.repo.organizationBySlug(slug)
}

func (s *Service) AllOrganizations(includeHidden bool) ([]Organization, error) {
	return s.repo.allOrganizations(includeHidden)
}

type passwordService interface {
	GenerateAndPersistPassword() (*security.CredentialSet, error)
	VerifyPassword(id int, password string) (bool, error)
	PersistCredentials(set security.PartialCredentialSet) (*security.CredentialSet, error)
}