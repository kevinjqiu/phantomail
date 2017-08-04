package smtpserver

import (
	"fmt"
	"net"
)

// SMTPServer represents an SMTP server
type SMTPServer struct {
	listenPort string
	listener   net.Listener
}

// Listen starts listening by creating a new listener
// This satisfies the TCPServer interface
func (s *SMTPServer) Listen() (net.Listener, error) {
	listener, err := net.Listen("tcp", s.listenPort)
	if err != nil {
		return nil, err
	}
	return listener, err
}

// Serve starts serving the requests using the listener
// This satisfies the TCPServer interface
func (s *SMTPServer) Serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func(c net.Conn) {
			fmt.Println("Got it")
			c.Close()
		}(conn)
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

// NewSMTPServer returns a new instance of SMTPServer type
func NewSMTPServer(cfg *Config) *SMTPServer {
	return &SMTPServer{
		listenPort: cfg.ListenPort,
	}
}
