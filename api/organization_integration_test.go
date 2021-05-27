package api

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

func (oits OrganizationIntegrationTestSuite) TestGetOrganization_shouldValidateAllowedReturnAttributes() {
	allowedAttrs := []string{"slug", "name", "description", "protected", "hidden", "tags" }
	// request a hidden organization since they have a larger attribute set
	hidden := oits.e.GET("/orgs/org-tests-hidden").WithHeader("plooral-credentials", "user-hidden").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object()

	// check that only allowed attributes should be returned
	hidden.Keys().Length().Equal(len(allowedAttrs))
	for _, attr := range allowedAttrs {
		hidden.ContainsKey(attr)
	}
}

func (oits OrganizationIntegrationTestSuite) TestOrganizationCRU() {
	// retrieve all organizations which are present
	oits.e.GET("/orgs").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(2)

	// create a new organization
	b, _ := ioutil.ReadFile("./test/data/request_bodies/valid_organization_creation_body.json")
	created := oits.e.POST("/orgs").
		WithBytes(b).
		Expect().
		Status(http.StatusCreated).JSON().Object().Value("data").Object()

	// use slug to look it up again and check org name
	slug := created.Value("slug").String().Raw()
	oits.e.GET("/orgs/" + slug).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("data").Object().Value("name").Equal("Integration Test Organization Creation")

	// update organization
	updateBody := created.Raw()
	updateBody["name"] = "Integration Test Organization Update"
	updated := oits.e.PATCH("/orgs/" + slug).
		WithJSON(updateBody).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Object()
	updated.Value("name").Equal("Integration Test Organization Update")

	// verify org list length
	oits.e.GET("/orgs").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(3)

	// verify GET on updated org
	oits.e.GET("/orgs/" + slug).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("data").Object().Value("name").Equal("Integration Test Organization Update")
}

func (oits OrganizationIntegrationTestSuite) TestProtectedOrganizationCRU() {
	// create a new organization
	b, _ := ioutil.ReadFile("./test/data/request_bodies/protected_organization_creation_body.json")
	created := oits.e.POST("/orgs").
		WithBytes(b).
		Expect().
		Status(http.StatusCreated).JSON().Object().Value("data").Object()

	created.Value("protected").Equal(true)
	created.Value("user_credentials").Object().Value("password").Equal("secret")
	created.Value("owner_credentials").Object().Value("password").NotNull()

	slug := created.Value("slug").String().Raw()
	oits.e.GET("/orgs/" + slug).
		Expect().
		Status(http.StatusForbidden).JSON().Object().Value("errors").Array().NotEmpty()

	oits.e.GET("/orgs/" + slug).WithHeader("plooral-credentials", "secret").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("slug").Equal(slug)
}

func (oits OrganizationIntegrationTestSuite) TestGetHiddenOrganization() {
	oits.e.GET("/orgs/org-tests-hidden").
		Expect().
		Status(http.StatusNotFound).JSON().Object().Value("errors").Array().NotEmpty()

	oits.e.GET("/orgs/org-tests-hidden").WithHeader("plooral-credentials", "user-hidden").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("slug").Equal("org-tests-hidden")
}

func (oits *OrganizationIntegrationTestSuite) SetupSuite() {
	oits.e = httpexpect.New(oits.Suite.T(), apiBaseURI())
}

func (oits *OrganizationIntegrationTestSuite) TearDownSuite() {
}

type OrganizationIntegrationTestSuite struct {
	suite.Suite
	e *httpexpect.Expect
}

func TestOrganizationIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationIntegrationTestSuite))
}