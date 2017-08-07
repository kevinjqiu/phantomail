package backends

import (
	"testing"

	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
)

func newTestController(input string) *caddy.Controller {
	c := caddy.NewTestController("smtp", input)
	c.Key = "smtp://localhost:2525"
	return c
}
func TestGetMaildirConfig___Rootpath(t *testing.T) {
	input := `maildir {
  rootPath "/tmp/maildir"
}`
	c := newTestController(input)
	maildirConfig, err := parseMaildirConfig(c)
	assert.Nil(t, err)
	assert.Equal(t, "/tmp/maildir", maildirConfig.rootPath)
}

func TestGetMaildirConfig___RootpathValueAbsent(t *testing.T) {
	input := `maildir {
  rootPath
}`
	c := newTestController(input)
	_, err := parseMaildirConfig(c)
	assert.NotNil(t, err)
}

func TestGetMaildirConfig___RootpathKeyAbsent(t *testing.T) {
	input := `maildir {
  someOtherKey
}`
	c := newTestController(input)
	_, err := parseMaildirConfig(c)
	assert.NotNil(t, err)
}
