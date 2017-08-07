package backends

import (
	"fmt"

	"github.com/kevinjqiu/phantomail/smtpserver"

	"github.com/mholt/caddy"
)

type stormConfig struct {
	path string
}

// NewStormStorageBackend creates a new instance of stormStorageBackend
func NewStormStorageBackend(c *caddy.Controller, next smtpserver.MessageHandler) StormStorageBackend {
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
	return StormStorageBackend{Next: next} // TODO
}

// StormStorageBackend stores incoming messages in a storm database
type StormStorageBackend struct {
	Next smtpserver.MessageHandler
}

// MessageReceived is the handler method of a `MessageHandler`
func (s StormStorageBackend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	fmt.Println("TODO: implement StormStorageBackend")
	return smtpserver.Next(s.Next, msg)
}

// Name is the name of the `MessageHandler`
func (s StormStorageBackend) Name() string {
	return "StormStorageBackend"
}
