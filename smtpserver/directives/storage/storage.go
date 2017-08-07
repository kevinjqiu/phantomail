package storage

import (
	"github.com/kevinjqiu/phantomail/smtpserver"
	"github.com/kevinjqiu/phantomail/smtpserver/directives/storage/backends"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("storage", caddy.Plugin{
		ServerType: smtpserver.ServerType,
		Action:     setupStorage,
	})
}

func parseStorageConfigs(c *caddy.Controller) error {
	config := smtpserver.GetConfig(c)

	for c.Next() {
		key := c.Val()
		if key != "storage" {
			return nil
		}
		args := c.RemainingArgs()
		if len(args) != 1 {
			return c.ArgErr()
		}
		storageType := args[0]
		if storageType == "storm" {
			config.AddMessageMiddleware(func(next smtpserver.MessageHandler) smtpserver.MessageHandler {
				return backends.NewStormStorageBackend(c, next)
			})
		} else if storageType == "maildir" {
		}
	}
	return nil
}

func setupStorage(c *caddy.Controller) error {
	err := parseStorageConfigs(c)
	if err != nil {
		return err
	}
	return nil
}
