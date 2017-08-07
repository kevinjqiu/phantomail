package logreceivedmessage

import (
	"log"

	"github.com/kevinjqiu/phantomail/smtpserver"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("logreceivedmessage", caddy.Plugin{
		ServerType: smtpserver.ServerType,
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	config := smtpserver.GetConfig(c)
	for c.Next() {
		key := c.Val()
		if key != "logreceivedmessage" {
			return nil
		}
		config.AddMessageMiddleware(func(next smtpserver.MessageHandler) smtpserver.MessageHandler {
			return logReceivedMessage{Next: next}
		})
	}
	return nil
}

type logReceivedMessage struct {
	Next smtpserver.MessageHandler
}

func (l logReceivedMessage) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	log.Printf("From: %s\n", msg.From)
	log.Printf("To: %s\n", msg.To)
	log.Printf("Received: %s\n", msg.Created)
	log.Printf("Content: %s\n", msg.Content.Body)
	return smtpserver.Next(l.Name(), l.Next, msg)
}

func (l logReceivedMessage) Name() string {
	return "LogReceivedMessage"
}
