package phantomail

import (
	// plugin the directives
	_ "github.com/kevinjqiu/phantomail/pkg/hostnames"
	_ "github.com/kevinjqiu/phantomail/pkg/logreceivedmessage"
	_ "github.com/kevinjqiu/phantomail/pkg/maildir"
	// plugin the server type
	_ "github.com/kevinjqiu/phantomail/pkg/smtpserver"
)
