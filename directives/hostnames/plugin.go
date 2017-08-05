package hostname

import (
	"github.com/kevinjqiu/phantomail/smtpserver"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("hostnames", caddy.Plugin{
		ServerType: "smtp",
		Action:     setupHost,
	})
}

func setupHost(c *caddy.Controller) error {
	config := smtpserver.GetConfig(c)
	if c.Key != "smtp" {
		return nil
	}

	for c.Next() {
		if !c.NextArg() {
			return c.ArgErr()
		}
		config.Hostnames = append(config.Hostnames, c.Val())
	}
	return nil
}
