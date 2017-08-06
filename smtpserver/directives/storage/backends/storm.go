package backends

import (
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

type StormStorageBackend struct {
	Next smtpserver.MessageHandler
}

func (s StormStorageBackend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	return smtpserver.NextOrFailure(s.Name(), s.Next, msg)
}

func (s StormStorageBackend) Name() string {
	return "StormStorageBackend"
}
