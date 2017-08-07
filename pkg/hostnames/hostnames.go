package hostnames

import (
	"github.com/kevinjqiu/phantomail/pkg/smtpserver"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("hostnames", caddy.Plugin{
		ServerType: smtpserver.ServerType,
		Action:     setupHost,
	})
}

func hostAlreadyExists(hostnames []string, newHost string) bool {
	for _, hostname := range hostnames {
		if hostname == newHost {
			return true
		}
	}
	return false
}

func setupHost(c *caddy.Controller) error {
	config := smtpserver.GetConfig(c)
	for c.Next() {
		if c.Val() != "hostnames" {
			return nil
		}
		args := c.RemainingArgs()
		if len(args) == 0 {
			return c.ArgErr()
		}

		for _, arg := range args {
			if hostAlreadyExists(config.Hostnames, arg) {
				continue
			}
			config.Hostnames = append(config.Hostnames, arg)
		}
	}
	return nil
}
