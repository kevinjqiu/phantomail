package smtpserver

import (
	"io"
	"log"
	"net"
	"strings"

	"github.com/mailhog/data"
	"github.com/mailhog/smtp"
	"github.com/mholt/caddy"
)

// SMTPServer represents an SMTP server
type SMTPServer struct {
	listener net.Listener
	config   *Config
}

// Listen starts listening by creating a new listener
// This satisfies the TCPServer interface
func (s *SMTPServer) Listen() (net.Listener, error) {
	listener, err := net.Listen("tcp", s.config.Bind())
	if err != nil {
		return nil, err
	}
	return listener, err
}

// Serve starts serving the requests using the listener
// This satisfies the TCPServer interface
func (s *SMTPServer) Serve(listener net.Listener) error {
	s.config.buildMiddlewareStacks()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		client := newSMTPClient(conn, s.config.rootMessageHandler)

		go client.handle()
	}
}

type smtpClient struct {
	conn               net.Conn
	writer             io.Writer
	reader             io.Reader
	line               string
	proto              *smtp.Protocol
	rootMessageHandler MessageHandler
}

func newSMTPClient(conn net.Conn, rootMessageHandler MessageHandler) *smtpClient {
	proto := smtp.NewProtocol()
	client := &smtpClient{
		conn:               conn,
		writer:             io.Writer(conn),
		reader:             io.Reader(conn),
		proto:              proto,
		rootMessageHandler: rootMessageHandler,
	}
	proto.MessageReceivedHandler = client.onMessageReceived
	// TODO: other handlers
	return client
}

// SMTPMessage represents an SMTP Message
type SMTPMessage struct {
	*data.SMTPMessage
}

func (client *smtpClient) onMessageReceived(msg *data.SMTPMessage) (id string, err error) {
	wrappedMessage := SMTPMessage{msg}
	return Next("[root]", client.rootMessageHandler, &wrappedMessage)
}

func (client *smtpClient) writeReply(reply *smtp.Reply) {
	for _, line := range reply.Lines() {
		client.writer.Write([]byte(line))
	}
}

func (client *smtpClient) read() bool {
	buf := make([]byte, 1024)
	n, err := client.reader.Read(buf)
	if n == 0 {
		// No more bytes to read
		return false
	}
	if err != nil {
		return false
	}

	text := string(buf[0:n])
	client.line += text

	for strings.Contains(client.line, "\r\n") {
		line, reply := client.proto.Parse(client.line)
		client.line = line

		if reply != nil {
			client.writeReply(reply)
			if reply.Status == 221 {
				client.conn.Close()
				return false
			}
		}
	}
	return true
}

func (client *smtpClient) handle() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Fatal error: %s\n", err)
		}
		client.conn.Close()
	}()
	reply := client.proto.Start()
	client.writeReply(reply)
	for client.read() == true {
	}
}

// ListenPacket listens and creates a UDP connection
// This satisfies the UDPServer interface
func (s *SMTPServer) ListenPacket() (net.PacketConn, error) {
	return nil, nil // Ignore
}

// ServePacket serves UDP requests
// This satisfies the UDPServer interface
func (s *SMTPServer) ServePacket(pc net.PacketConn) error {
	return nil // Ignore
}

// Stop closes the listening socket
func (s *SMTPServer) Stop() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

// OnStartupComplete shows the current effective configuration
func (s *SMTPServer) OnStartupComplete() {
	if !caddy.Quiet {
		log.Printf("SMTP server started with configuration: %v\n", s.config)
	}
}

// NewSMTPServer returns a new instance of SMTPServer type
func NewSMTPServer(cfg *Config) *SMTPServer {
	cfg.buildMiddlewareStacks()
	return &SMTPServer{
		config: cfg,
	}
}

// GetConfig gets the Config given the controller
func GetConfig(c *caddy.Controller) *Config {
	ctx := c.Context().(*smtpContext)
	key := c.Key
	if c.ServerType() == ServerType {
		if cfg, ok := ctx.keysToConfigs[key]; ok {
			return cfg
		}
	}
	// test config
	testConfig := &Config{}
	ctx.keysToConfigs[key] = testConfig
	return testConfig
}
