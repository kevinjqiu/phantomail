package smtpserver

import (
	"fmt"
	"log"

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

func (c *smtpContext) InspectServerBlocks(sourceFile string, serverBlocks []caddyfile.ServerBlock) ([]caddyfile.ServerBlock, error) {
	log.Println(serverBlocks)
	currentKey := ""
	cfg := make(map[string][]string)
	for _, serverBlock := range serverBlocks {
		for _, key := range serverBlock.Keys {
			if _, dup := c.keysToConfigs[key]; dup {
				return serverBlocks, fmt.Errorf("duplicate key: %s", key)
			}

			switch key {
			case "smtp":
				currentKey = key
				cfg[currentKey] = []string{}
			default:
				cfg[currentKey] = append(cfg[currentKey], key)
			}
		}
	}

	for k, v := range cfg {
		if len(v) == 0 {
			return serverBlocks, fmt.Errorf("invalid configuration: %s", k)
		}

		smtpServerConfig := &Config{
			ListenPort: v[0],
		}

		c.saveConfig(k, smtpServerConfig)
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
}

func newContext() caddy.Context {
	return &smtpContext{keysToConfigs: make(map[string]*Config)}
}
