package publisher

import (
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
)

type BasePublisher struct {
	Logger *logger.Logger
	url    string
	config *config.Config
}
