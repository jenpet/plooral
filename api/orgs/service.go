package orgs

import (
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"net/http"
)

type Service struct {
	repo *repository
}

func newDefaultService() *Service {
	return &Service{repo: newDefaultRepository()}
}

func newService(r *repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateOrganization(o partialOrganization) (*Organization, error) {
	org, err := s.OrganizationBySlug(*o.Slug)
	if err != nil && errors.ErrKind(err) != database.KNoEntityFound {
		return nil, err
	}
	if org != nil {
		return nil, errors.Ef("Organization '%s' does already exist.", http.StatusBadRequest, rest.KUserInputInvalid)
	}
	return s.repo.upsertOrganization(o.toOrganization())
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
