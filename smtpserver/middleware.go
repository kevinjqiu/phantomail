package smtpserver

import "fmt"

type (
	// MessageMiddleware is the middleware for handling received messages
	MessageMiddleware func(MessageHandler) MessageHandler
	// MessageHandler handles received messages
	MessageHandler interface {
		MessageReceived(*SMTPMessage) (string, error)
		Name() string
	}
	// MessageHandlerFunc is a convenience type for defining message handlers
	MessageHandlerFunc func(*SMTPMessage) (string, error)
)

// NextOrFailure chains the middlewares together and executes them
func NextOrFailure(name string, next MessageHandler, msg *SMTPMessage) (id string, err error) {
	if next != nil {
		return next.MessageReceived(msg)
	}
	return "", fmt.Errorf("No handlers available: %s", name)
}
