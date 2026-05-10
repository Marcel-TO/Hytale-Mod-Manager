package publisher

import (
	"marcel-to/hytale/mod-publisher/config"
	"marcel-to/hytale/mod-publisher/logger"
)

type BasePublisher struct {
	Logger logger.Logger
	url    string
	config *config.Config
}
