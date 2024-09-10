package config

import (
	"time"
)

const (
	WebsocketTimeout   = time.Second * 600
	WebsocketKeepAlive = time.Second * 86400
	WebsocketMaxBytes  = 655350
)
