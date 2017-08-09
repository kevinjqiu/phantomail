package backends

import (
	"fmt"
	"log"

	"github.com/kevinjqiu/phantomail/pkg/smtpserver"
	"github.com/mholt/caddy"
	md "github.com/sloonz/go-maildir"
	"os"
	"strconv"
	"os/user"
)

func init() {
	caddy.RegisterPlugin("maildir", caddy.Plugin{
		ServerType: smtpserver.ServerType,
		Action:     setup,
	})
}

type config struct {
	rootPath string
	mode     os.FileMode
	uid      int
	gid      int
}

func (c config) validate() error {
	if c.rootPath == "" {
		return fmt.Errorf("`rootPath` is mandatory for `maildir`")
	}
	return nil
}

type backend struct {
	Next   smtpserver.MessageHandler
	Config config
}

// Name is the name of the `MessageHandler`
func (m backend) Name() string {
	return "MaildirStorageBackend"
}

// MessageReceived is the handler method of a `MessageHandler`
func (m backend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	for _, recipient := range msg.To {
		maildirPath := fmt.Sprintf("%s/%s@%s", m.Config.rootPath, recipient.Mailbox, recipient.Domain)
		maildir, err := md.NewWithPerm(maildirPath, true, m.Config.mode, m.Config.uid, m.Config.gid)
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
		return backend{Next: next, Config: maildirConfig}
	})
	return nil
}

func parseMaildirConfig(c *caddy.Controller) (config, error) {
	config := config{
		mode: os.FileMode(0755),
		uid: os.Getuid(),
		gid: os.Getgid(),
	}
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
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				config.rootPath = args[0]
			case "mode":
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				mode, err := strconv.ParseInt(args[0], 8, 32)
				if err != nil {
					return config, err
				}
				config.mode = os.FileMode(mode)
			case "user":
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				username := args[0]
				u, err := user.Lookup(username)
				if err != nil {
					return config, fmt.Errorf("Cannot find user %v", username)
				}
				uid, err := strconv.Atoi(u.Uid)
				if err != nil {
					return config, fmt.Errorf("uid is not an integer. Are you on POSIX?")
				}
				config.uid = uid
			case "group":
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				groupname := args[0]
				g, err := user.LookupGroup(groupname)
				if err != nil {
					return config, fmt.Errorf("Cannot find group %v", groupname)
				}
				gid, err := strconv.Atoi(g.Gid)
				if err != nil {
					return config, fmt.Errorf("gid is not an integer. Are you on POSIX?")
				}
				config.gid = gid
			}
		}
	}

	err := config.validate()
	return config, err
}
