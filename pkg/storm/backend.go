package storm

import (
	"fmt"

	"github.com/kevinjqiu/phantomail/pkg/smtpserver"
)

type stormConfig struct {
	path string
}

// StormStorageBackend stores incoming messages in a storm database
type stormStorageBackend struct {
	Next   smtpserver.MessageHandler
	config stormConfig
}

// MessageReceived is the handler method of a `MessageHandler`
func (s stormStorageBackend) MessageReceived(msg *smtpserver.SMTPMessage) (string, error) {
	fmt.Println("TODO: implement StormStorageBackend")
	return smtpserver.Next(s.Next, msg)
}

// Name is the name of the `MessageHandler`
func (s stormStorageBackend) Name() string {
	return "StormStorageBackend"
}
