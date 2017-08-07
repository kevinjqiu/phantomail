package backends

import (
	"fmt"

	"github.com/kevinjqiu/phantomail/smtpserver"
	"github.com/mholt/caddy"
)

type maildirConfig struct {
	rootPath string
}

// MaildirStorageBackend is the maildir storage backend
type MaildirStorageBackend struct {
	Next   smtpserver.MessageHandler
	config maildirConfig
}

// Name is the name of the `MessageHandler`
func (m MaildirStorageBackend) Name() string {
	return "MaildirStorageBackend"
}

// MessageReceived is the handler method of a `MessageHandler`
func (m MaildirStorageBackend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	return smtpserver.Next(m.Next, msg)
}

func getMaildirConfig(c *caddy.Controller) (maildirConfig, error) {
	config := maildirConfig{}
	for c.NextBlock() {
		key := c.Val()
		args := c.RemainingArgs()
		switch key {
		case "rootPath":
			if len(args) == 1 {
				config.rootPath = args[0]
			} else {
				return config, c.ArgErr()
			}
		}
	}

	if config.rootPath == "" {
		return config, fmt.Errorf("`rootPath` is mandatory for `maildir`")
	}
	return config, nil
}

// NewMaildirBackend creates a new instance of `MaildirStorageBackend`
func NewMaildirBackend(c *caddy.Controller, next smtpserver.MessageHandler) (MaildirStorageBackend, error) {
	config, err := getMaildirConfig(c)
	if err != nil {
		return MaildirStorageBackend{}, nil
	}
	return MaildirStorageBackend{Next: next, config: config}, nil
}
