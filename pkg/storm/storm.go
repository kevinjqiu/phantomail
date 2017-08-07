package storm

import (
	"github.com/kevinjqiu/phantomail/pkg/smtpserver"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("storm", caddy.Plugin{
		ServerType: smtpserver.ServerType,
		Action:     setup,
	})
}

func parseStormConfig(c *caddy.Controller) (stormConfig, error) {
	config := stormConfig{}
	for c.Next() {
		key := c.Val()
		if key != "storm" {
			return config, nil
		}
		for c.NextBlock() {
			key := c.Val()
			args := c.RemainingArgs()
			switch key {
			case "path":
				if len(args) == 1 {
					config.path = args[0]
				} else {
					return config, c.ArgErr()
				}
			}
		}
	}
	// TODO make sure mandatory args are supplied
	return config, nil
}

func setup(c *caddy.Controller) error {
	stormCfg, err := parseStormConfig(c)
	if err != nil {
		return err
	}
	config := smtpserver.GetConfig(c)
	config.AddMessageMiddleware(func(next smtpserver.MessageHandler) smtpserver.MessageHandler {
		return stormStorageBackend{Next: next, config: stormCfg}
	})
	return nil
}
