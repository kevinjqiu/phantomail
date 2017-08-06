package smtpserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

const serverType = "smtp"

var directives = []string{"hostnames"}

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

func (c *smtpContext) saveConfig(key string, cfg *Config) {
	c.configs = append(c.configs, cfg)
	c.keysToConfigs[key] = cfg
}

const defaultSMTPPort = "25"

func parseSMTPLabel(label string) (string, string, error) {
	label = strings.ToLower(label)
	if !strings.HasPrefix(label, "smtp://") {
		return "", "", fmt.Errorf("%s does not start with 'smtp://'", label)
	}
	label = strings.TrimPrefix(label, "smtp://")
	parts := strings.Split(label, ":")
	if len(parts) == 0 || len(parts) > 2 {
		return "", "", fmt.Errorf("%s does not contain a binding IP and port", label)
	}
	bindAddr := parts[0]
	if bindAddr == "*" {
		bindAddr = "0.0.0.0"
	}
	if len(parts) == 1 {
		return bindAddr, defaultSMTPPort, nil
	}
	port := parts[1]
	_, err := strconv.Atoi(port)
	if err != nil {
		return "", "", err
	}
	return bindAddr, port, nil
}

func (c *smtpContext) InspectServerBlocks(sourceFile string, serverBlocks []caddyfile.ServerBlock) ([]caddyfile.ServerBlock, error) {
	for _, serverBlock := range serverBlocks {
		for _, key := range serverBlock.Keys {
			if _, dup := c.keysToConfigs[key]; dup {
				return serverBlocks, fmt.Errorf("duplicate key: %s", key)
			}

			bindAddr, port, err := parseSMTPLabel(key)
			if err != nil {
				return serverBlocks, err
			}
			smtpServerConfig := &Config{
				ListenPort: port,
				BindAddr:   bindAddr,
			}
			c.saveConfig(key, smtpServerConfig)
		}
	}
	return serverBlocks, nil
}

func (c *smtpContext) MakeServers() ([]caddy.Server, error) {
	var servers []caddy.Server

	for _, cfg := range c.configs {
		servers = append(servers, NewSMTPServer(cfg))
	}
	return servers, nil
}

// Config contains configuration details about an SMTP server type
type Config struct {
	Hostnames  []string
	ListenPort string
	BindAddr   string
}

func newContext() caddy.Context {
	return &smtpContext{keysToConfigs: make(map[string]*Config)}
}
