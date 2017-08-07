package backends

import (
	"fmt"
	"log"

	"github.com/kevinjqiu/phantomail/pkg/smtpserver"
	"github.com/mholt/caddy"
	md "github.com/sloonz/go-maildir"
)

func init() {
	caddy.RegisterPlugin("maildir", caddy.Plugin{
		ServerType: smtpserver.ServerType,
		Action:     setup,
	})
}

type maildirConfig struct {
	rootPath string
}

// MaildirStorageBackend is the maildir storage backend
type maildirStorageBackend struct {
	Next   smtpserver.MessageHandler
	Config maildirConfig
}

// Name is the name of the `MessageHandler`
func (m maildirStorageBackend) Name() string {
	return "MaildirStorageBackend"
}

// MessageReceived is the handler method of a `MessageHandler`
func (m maildirStorageBackend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	for _, recipient := range msg.To {
		maildirPath := fmt.Sprintf("%s/%s@%s", m.Config.rootPath, recipient.Mailbox, recipient.Domain)
		maildir, err := md.New(maildirPath, true)
		if err != nil {
			return "", err
		}
		filename, err := maildir.CreateMail(msg.Bytes())
		if err != nil {
			return "", err
		}
		log.Printf("Message saved to %s\n", filename)
	}
	return smtpserver.Next(m.Next, msg)
}

func setup(c *caddy.Controller) error {
	maildirConfig, err := parseMaildirConfig(c)
	if err != nil {
		return err
	}
	config := smtpserver.GetConfig(c)
	config.AddMessageMiddleware(func(next smtpserver.MessageHandler) smtpserver.MessageHandler {
		return maildirStorageBackend{Next: next, Config: maildirConfig}
	})
	return nil
}

func parseMaildirConfig(c *caddy.Controller) (maildirConfig, error) {
	config := maildirConfig{}
	for c.Next() {
		key := c.Val()
		if key != "maildir" {
			return config, nil
		}

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
	}

	if config.rootPath == "" {
		return config, fmt.Errorf("`rootPath` is mandatory for `maildir`")
	}
	return config, nil
}
