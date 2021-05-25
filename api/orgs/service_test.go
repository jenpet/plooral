package orgs

import (
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"github.com/jenpet/plooral/test/data"
	"github.com/jenpet/plooral/testutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

func (ost OrganizationServiceTestSuite) TestAllOrganizations_shouldReturnFullList() {
	// request without hidden
	orgs, err := ost.cut.AllOrganizations(false)
	ost.NoError(err, "no error expected during lookup")
	ost.Len(orgs, 2, "expected two organizations which are not hidden")
	for _, o := range orgs {
		ost.NotEqual("org-tests-hidden", o.Slug, "expected no hidden org in result set")
	}

	// request with hidden
	orgs, err = ost.cut.AllOrganizations(true)
	ost.NoError(err, "no error expected during lookup")
	ost.Len(orgs, 3, "expected two organizations which are not hidden")
	foundHidden := false
	for _, o := range orgs {
		if o.Slug == "org-tests-hidden" {
			foundHidden = true
		}
	}
	ost.True(foundHidden, "expected hidden organization in result set")
}

func (ost OrganizationServiceTestSuite) TestOrganizationBySlug_shouldReturnOrganization() {
	o, err := ost.cut.OrganizationBySlug("org-tests-regular")
	ost.NoError(err, "no error expected during lookup")
	ost.Equal(o.Slug, "org-tests-regular")
	ost.False(o.Protected)
	ost.False(o.Hidden)
	o = nil
	err = nil
	o, err = ost.cut.OrganizationBySlug("unknown")
	ost.Equal(database.KNoEntityFound, errors.ErrKind(err), "common database error expected")
	ost.Nil(o, "expected returned organization to be nil")
}

func (ost OrganizationServiceTestSuite) TestCreateOrganization_shouldInsertAndReturnOrgWithNewID() {
	o := partialOrganization{}
	o.setSlug("org-tests-insert")
	o.setName("Inserted Organization")
	o.setDescription("")
	o.setHidden(false)
	o.setProtected(false)
	o.setTags([]string{})
	// insert organization
	inserted, err := ost.cut.CreateOrganization(o)
	ost.NoError(err, "no error expected")
	ost.True(inserted.ID >= 0, "id should be set")

	// update organization has to fail
	recreated, err := ost.cut.CreateOrganization(o)
	ost.Error(err, "error expected")
	ost.Equal(rest.KUserInputInvalid, errors.ErrKind(err))
	ost.Nil(recreated, "expected recreated to be nil")

	lookup, err := ost.cut.OrganizationBySlug("org-tests-insert")
	ost.NoError(err, "no error expected")
	ost.NotNil(lookup, "expected lookup not to be nil")
}

func (ost OrganizationServiceTestSuite) TestUpdateOrganization_shouldUpdateAndReturnOrg() {
	o := partialOrganization{}
	o.setSlug("non-existent")
	_, err := ost.cut.UpdateOrganization(o)
	ost.Error(err, "error expected")

	o = partialOrganization{}
	o.setName("Updated Title")
	o.setSlug("org-tests-regular")

	updated, err := ost.cut.UpdateOrganization(o)
	ost.NoError(err, "no error expected")
	ost.NotNil(updated, "expected result not to be nil")

	lookup, _ := ost.cut.OrganizationBySlug("org-tests-regular")
	ost.Equal("Updated Title", lookup.Name)
}

type OrganizationServiceTestSuite struct {
	suite.Suite
	cut *Service
	repo *repository
	pgName string
}

func (ost *OrganizationServiceTestSuite) SetupSuite() {
	pgName, pgURI := testutil.NewPostgres()
	ost.pgName = pgName
	// apply the default migrations
	database.ApplyDefaultMigrations(log.StandardLogger(), pgURI)
	testutil.ApplyFromEmbeddedFS(pgURI, data.TestDataMigrations, "migrations")
	ost.repo = newRepository(pgURI)
	ost.cut = newService(ost.repo)
}

func (ost *OrganizationServiceTestSuite) TearDownSuite() {
	testutil.Purge(ost.pgName)
}

func TestOrganizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationServiceTestSuite))
}

