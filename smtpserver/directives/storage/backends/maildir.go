package backends

import (
	"fmt"
	"log"

	"github.com/kevinjqiu/phantomail/smtpserver"
	"github.com/mholt/caddy"
	md "github.com/sloonz/go-maildir"
)

// MaildirConfig is the configuration for the maildir storage backend
type MaildirConfig struct {
	rootPath string
}

// MaildirStorageBackend is the maildir storage backend
type MaildirStorageBackend struct {
	Next   smtpserver.MessageHandler
	Config MaildirConfig
}

// Name is the name of the `MessageHandler`
func (m MaildirStorageBackend) Name() string {
	return "MaildirStorageBackend"
}

// MessageReceived is the handler method of a `MessageHandler`
func (m MaildirStorageBackend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	maildir, err := md.New(m.Config.rootPath, true)
	if err != nil {
		return "", err
	}
	filename, err := maildir.CreateMail(msg.Bytes())
	if err != nil {
		return "", err
	}
	log.Printf("Saved message to %s\n", filename)
	return smtpserver.Next(m.Next, msg)
}

// ParseMaildirConfig parses the maildir configuration
func ParseMaildirConfig(c *caddy.Controller) (MaildirConfig, error) {
	config := MaildirConfig{}
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
