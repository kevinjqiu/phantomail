package hostnames

import (
	"testing"

	"github.com/kevinjqiu/phantomail/pkg/smtpserver"
	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
)

func newTestController(input string) *caddy.Controller {
	c := caddy.NewTestController("smtp", input)
	c.Key = "smtp"
	return c
}
func TestSetupHost___NoHostnameAfterDirective(t *testing.T) {
	c := newTestController(`hostnames`)
	err := setupHost(c)
	assert.NotNil(t, err)
}

func TestSetupHost___SingleHostname(t *testing.T) {
	c := newTestController(`hostnames "phantomail.com"`)
	err := setupHost(c)
	assert.Nil(t, err)
	cfg := smtpserver.GetConfig(c)
	assert.EqualValues(t, []string{"phantomail.com"}, cfg.Hostnames)
}

func TestSetupHost___MultipleHostnames(t *testing.T) {
	c := newTestController(`hostnames "phantomail.com" "m.phantomail.com"`)
	err := setupHost(c)
	assert.Nil(t, err)
	cfg := smtpserver.GetConfig(c)
	assert.EqualValues(t, []string{"phantomail.com", "m.phantomail.com"}, cfg.Hostnames)
}
