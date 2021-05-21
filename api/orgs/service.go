package orgs

type Service struct {
	repo *repository
}

func newDefaultService() *Service {
	return &Service{repo: newDefaultRepository()}
}

func newService(r *repository) *Service {
	return &Service{repo: r}
}

func (s *Service) UpsertOrganization(o Organization) (*Organization, error) {
	return s.repo.upsertOrganization(o)
}

func (s *Service) OrganizationBySlug(slug string) (*Organization, error) {
	return s.repo.organizationBySlug(slug)
}

func (s *Service) AllOrganizations(includeHidden bool) ([]Organization, error) {
	return s.repo.allOrganizations(includeHidden)
}
