package storage

import (
	"github.com/kevinjqiu/phantomail/smtpserver/directives/storage/backends"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("storage", caddy.Plugin{
		ServerType: "smtp",
		Action:     setupStorage,
	})
}

// Config is the interface for configs for different storage backends
type Config interface {
}

func parseStorageConfigs(c *caddy.Controller) (*[]Config, error) {
	// config := smtpserver.GetConfig(c)

	var storageConfigs []Config
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return nil, c.ArgErr()
		}
		storageType := args[0]
		if storageType == "storm" {
			c, err := backends.NewStormConfig(c)
			if err != nil {
				return nil, err
			}
			storageConfigs = append(storageConfigs, c)
		} else if storageType == "maildir" {
		}
	}
	return &storageConfigs, nil
}

func setupStorage(c *caddy.Controller) error {
	_, err := parseStorageConfigs(c)
	if err != nil {
		return err
	}
	return nil
}
