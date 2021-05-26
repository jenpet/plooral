package security

import (
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/test/data"
	"github.com/jenpet/plooral/testutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

func (sst *SecurityServiceTestSuite) TestVerifyPassword_shouldOnlyConsiderActiveOnesAsValid() {
	ok, err := sst.cut.VerifyPassword(1, "active-password")
	sst.Nil(err, "no error expected")
	sst.True(ok, "expected active password to be verified successfully")

	ok, err = sst.cut.VerifyPassword(2, "inactive-password")
	sst.Nil(err, "no error expected")
	sst.False(ok, "expected inactive password to not be verified successfully")

	ok, err = sst.cut.VerifyPassword(3, "nonexistent-password")
	sst.Nil(err, "no error expected")
	sst.False(ok, "expected nonexistent password to not be verified successfully")
}

func (sst *SecurityServiceTestSuite) TestPersistPassword_shouldPersistAndReturnSet() {
	pw := "password"
	pw2 := "password2"
	set, err := sst.cut.PersistCredentials(PartialCredentialSet{
		Password:             &pw,
		PasswordConfirmation: &pw2,
	})
	sst.Nil(set, "set expected to be nil when passwords are unequal")
	sst.Error(err, "error expected when passwords are unequal")

	set, err = sst.cut.PersistCredentials(PartialCredentialSet{
		Password:             &pw,
		PasswordConfirmation: &pw,
	})
	sst.NoError(err, "no error expected when passwords are equal")
	sst.Equal(pw, set.Password, "expected returned password to match input")
	sst.True(set.ID > 0, "expected returned ID to have a valid value")

	ok, _ := sst.cut.VerifyPassword(set.ID, set.Password)
	sst.True(ok, "expected persisted password to be present")
}

func (sst SecurityServiceTestSuite) TestGenerateAndPersistPassword_shouldReturnSet() {
	set, err := sst.cut.GenerateAndPersistPassword()
	sst.NoError(err, "no error expected when generated")
	sst.NotNil(set, "expected returned set not to be nil")

	ok, _ := sst.cut.VerifyPassword(set.ID, set.Password)
	sst.True(ok, "expected persisted password to be present")
}

func (sst *SecurityServiceTestSuite) SetupSuite() {
	pgName, pgURI := testutil.NewPostgres()
	sst.pgName = pgName
	// apply the default migrations
	database.ApplyDefaultMigrations(log.StandardLogger(), pgURI)
	testutil.ApplyFromEmbeddedFS(pgURI, data.TestDataMigrations, "migrations")
	sst.repo = newRepository(pgURI)
	sst.cut = newService(sst.repo)
}

func (sst *SecurityServiceTestSuite) TearDownSuite() {
	testutil.Purge(sst.pgName)
}

type SecurityServiceTestSuite struct {
	suite.Suite
	cut *Service
	repo *repository
	pgName string
}

func TestSecurityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityServiceTestSuite))
}
