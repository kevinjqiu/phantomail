package phantomail

import (
	// plugin the directives
	_ "github.com/kevinjqiu/phantomail/directives/hostnames"
	// plugin the server type
	_ "github.com/kevinjqiu/phantomail/smtpserver"
)
