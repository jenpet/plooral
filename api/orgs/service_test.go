package orgs

import (
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"github.com/jenpet/plooral/security"
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
	ost.Nil(inserted.UserSecurity, "no user password expected when org is not hidden or protected")
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

func (ost OrganizationServiceTestSuite) TestCreateOrganization_whenProtectedEnabled_shouldReturnOrgWithSecurity() {
	o := partialOrganization{}
	o.setSlug("org-tests-pw-insert")
	o.setName("Inserted Organization")
	o.setDescription("")
	o.setHidden(true)
	o.setProtected(false)
	o.setTags([]string{})

	// hidden but not protected orgs should result in an error
	inserted, err := ost.cut.CreateOrganization(o)
	ost.Nil(inserted, "expected org to be nil")
	ost.Equal(rest.KUserInputInvalid, errors.ErrKind(err), "error expected when org is hidden but not protected")

	// protected but no user password
	o.setProtected(true)
	inserted, err = ost.cut.CreateOrganization(o)
	ost.Nil(inserted, "expected org to be nil")
	ost.Equal(rest.KUserInputInvalid, errors.ErrKind(err), "error expected when org is hidden but not protected")

	o.setPassword("pw")
	o.setPasswordConfirmation("pw")
	inserted, err = ost.cut.CreateOrganization(o)
	ost.NotNil(inserted, "expected inserted org not to be nil")
	ost.NotNil(inserted.UserSecurity, "expected user security to be set after creation")
	ost.NotNil(inserted.OwnerSecurity, "expected owner security to be set after creation")
	ost.Equal(inserted.UserSecurity.Password, "pw", "given password should be set for user security")
	ost.Equal(inserted.OwnerSecurity.Password, "generated", "generated password should be set for owner security")
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

type mockedPasswordService struct {}

func (mps *mockedPasswordService) GenerateAndPersistPassword() (*security.CredentialSet, error) {
	return &security.CredentialSet{ ID: 1, Password: "generated"}, nil
}

func (mps *mockedPasswordService) PersistCredentials(set security.PartialCredentialSet) (*security.CredentialSet, error) {
	return &security.CredentialSet{ ID: 2, Password: *set.Password }, nil
}

func (mps *mockedPasswordService) VerifyPassword(id int, password string) (bool, error) {
	if password == "verified" {
		return true, nil
	}
	return false, nil
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
	ost.cut = newService(ost.repo, &mockedPasswordService{})
}

func (ost *OrganizationServiceTestSuite) TearDownSuite() {
	testutil.Purge(ost.pgName)
}

func TestOrganizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationServiceTestSuite))
}

