package api

import (
	"fmt"
	"github.com/jenpet/plooral/test/data"
	"github.com/jenpet/plooral/testutil"
	"net/http/httptest"
	"os"
	"testing"
)

var apiURI string
var pgURI string

func TestMain(m *testing.M) {
	// startup postgres db
	var pgName string
	pgName, pgURI = testutil.NewPostgres()

	// set db var separately since it will be generated on runtime
	testutil.SetEnvs(map[string]string{"POSTGRES_URI": pgURI})

	srv := httptest.NewServer(Server())
	defer srv.Close()

	// set URI for requests
	apiURI = srv.URL

	// seed dummy data
	testutil.ApplyFromEmbeddedFS(pgURI, data.TestDataMigrations, "migrations")

	code := m.Run()

	testutil.Purge(pgName)
	os.Exit(code)
}

// base URI of the API based on the actual test server
// Used to perform actual calls targeting the API
func apiBaseURI() string {
	return fmt.Sprintf("%s/%s", apiURI, "api/v1")
}