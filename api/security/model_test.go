package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrganization_IsValid(t *testing.T) {
	var csp *CredentialSet
	assert.False(t, csp.IsValid())
	cs := CredentialSet{ ID: 0, Password: ""}
	assert.False(t, cs.IsValid())
	cs.ID = 1
	cs.Password = "pw"
	assert.True(t, cs.IsValid())
}
