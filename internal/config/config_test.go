package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const envFile = "../../.env"

func TestConfig(t *testing.T) {
	want := Config{
		Port: "8080",
		LogLevel: "info",
	}

	got, err := Get(envFile)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, want, got, "Configuration file (.env) isn't exist or didn't find")
}
