package storage

import (
	"fmt"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("storage", caddy.Plugin{
		ServerType: "smtp",
		Action:     setupStorage,
	})
}
func setupStorage(c *caddy.Controller) error {
	fmt.Println(c)
	return nil
}
