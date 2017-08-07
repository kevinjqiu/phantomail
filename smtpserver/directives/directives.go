package directives

import (
	// Plugins
	_ "github.com/kevinjqiu/phantomail/smtpserver/directives/hostnames"
	_ "github.com/kevinjqiu/phantomail/smtpserver/directives/logreceivedmessage"
	_ "github.com/kevinjqiu/phantomail/smtpserver/directives/storage"
)
