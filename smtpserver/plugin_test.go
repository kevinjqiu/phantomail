package smtpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSMTPLabel___NotSMTP(t *testing.T) {
	_, _, err := parseSMTPLabel("http://localhost:8080")
	assert.NotNil(t, err)
}

func TestParseSMTPLabel___NoPort(t *testing.T) {
	host, port, err := parseSMTPLabel("smtp://localhost")
	assert.Nil(t, err)
	assert.Equal(t, "localhost", host)
	assert.Equal(t, "25", port)
}

func TestParseSMTPLabel___WildcardHost(t *testing.T) {
	host, port, err := parseSMTPLabel("smtp://*")
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", host)
	assert.Equal(t, "25", port)
}

func TestParseSMTPLabel___WildcardHostWithPort(t *testing.T) {
	host, port, err := parseSMTPLabel("smtp://*:2525")
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", host)
	assert.Equal(t, "2525", port)
}
