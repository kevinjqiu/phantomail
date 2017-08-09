package backends

import (
	"testing"

	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
	"os"
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

func TestGetMaildirConfig___DefaultModeUidGid(t *testing.T) {
	input := `maildir {
  rootPath "/tmp/maildir"
}`
	c := newTestController(input)
	config, err := parseMaildirConfig(c)
	assert.Nil(t, err)
	assert.Equal(t, os.FileMode(0755), config.mode)
	assert.Equal(t, os.Getuid(), config.uid)
	assert.Equal(t, os.Getgid(), config.gid)
}

func TestGetMaildirConfig___CustomMode(t *testing.T) {
	input := `maildir {
  rootPath "/tmp/maildir"
  mode 0777
}`
	c := newTestController(input)
	config, err := parseMaildirConfig(c)
	assert.Nil(t, err)
	assert.Equal(t, os.FileMode(0777), config.mode)
	assert.Equal(t, os.Getuid(), config.uid)
	assert.Equal(t, os.Getgid(), config.gid)
}

func TestGetMaildirConfig___CustomUidGid(t *testing.T) {
	input := `maildir {
  rootPath "/tmp/maildir"
  user "root"
  group "root"
}`
	c := newTestController(input)
	config, err := parseMaildirConfig(c)
	assert.Nil(t, err)
	assert.Equal(t, 0, config.uid)
	assert.Equal(t, 0, config.gid)
}
