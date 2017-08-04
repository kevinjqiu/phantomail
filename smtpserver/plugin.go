package smtpserver

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

const serverType = "smtp"

var directives = []string{"host"}

func init() {
	caddy.RegisterServerType(serverType, caddy.ServerType{
		Directives: func() []string { return directives },
		DefaultInput: func() caddy.Input {
			return caddy.CaddyfileInput{
				ServerTypeName: serverType,
			}
		},
		NewContext: newContext,
	})
}

type smtpContext struct {
	keysToConfigs map[string]*Config
	configs       []*Config
}

func (c *smtpContext) InspectServerBlocks(sourceFile string, serverBlocks []caddyfile.ServerBlock) ([]caddyfile.ServerBlock, error) {
	return serverBlocks, nil // TODO
}

func (c *smtpContext) MakeServers() ([]caddy.Server, error) {
	var servers []caddy.Server

	servers = append(servers, NewSMTPServer())
	return servers, nil
}

// Config contains configuration details about an SMTP server type
type Config struct {
	Hostname   string
	ListenPort string
}

func newContext() caddy.Context {
	return &smtpContext{keysToConfigs: make(map[string]*Config)}
}
