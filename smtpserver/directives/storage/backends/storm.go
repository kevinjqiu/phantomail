package backends

import (
	"fmt"

	"github.com/mholt/caddy"
)

type stormConfig struct {
	path string
}

func NewStormConfig(c *caddy.Controller) (*stormConfig, error) {
	config := stormConfig{}
	for c.NextBlock() {
		key := c.Val()
		args := c.RemainingArgs()

		switch key {
		case "path":
			if len(args) == 1 {
				config.path = args[0]
			}
		}
	}
	fmt.Println(config)
	return &config, nil
}

type StormBackend struct {
}
