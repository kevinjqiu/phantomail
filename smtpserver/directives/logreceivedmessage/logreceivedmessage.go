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
	// m := msg.Parse(client.proto.Hostname)
	m := msg.Parse("")
	log.Printf("From: %s\n", m.From)
	log.Printf("To: %s\n", m.To)
	log.Printf("Received: %s\n", m.Created)
	log.Printf("Content: %s\n", m.Content.Body)
	return smtpserver.Next(l.Name(), l.Next, msg)
}

func (l logReceivedMessage) Name() string {
	return "LogReceivedMessage"
}
