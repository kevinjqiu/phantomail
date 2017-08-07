package smtpserver

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

// Next chains the middlewares together and executes them
func Next(next MessageHandler, msg *SMTPMessage) (string, error) {
	if next != nil {
		return next.MessageReceived(msg)
	}
	return "", nil
}
