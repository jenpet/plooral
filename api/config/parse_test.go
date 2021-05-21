package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParse_whenSpecIsNotAPtr_shouldPanic(t *testing.T) {
	assert.Panics(t, func() {Parse("")}, "parse should panic when not passing a pointer")
	assert.Panics(t, func() {Parse(nil)}, "parse should panic when passing nil")
}

func TestParse_shouldReadFromEnvVars(t *testing.T) {
	_ = os.Setenv("TEST_VAR", "true")
	cfg := testConfig{}
	Parse(&cfg)
	assert.Equal(t, true, cfg.TestVar, "parse should parse passed config from env vars")
}

type testConfig struct {
	TestVar bool `split_words:"true"`
}
