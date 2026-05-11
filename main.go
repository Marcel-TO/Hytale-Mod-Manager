package main

import (
	"fmt"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/publisher"
)

func main() {
	logger := logger.NewLogger("HMP", true)
	logger.Info("Hytale Mod Publisher")

	cfg := config.Config{}
	cfg.GetConfig(*logger, "publish-config.yaml")
	logger.Info("Configuration loaded successfully!")

	logger.Info("Publishing to CurseForge...")
	logger.Info(fmt.Sprintf("Found %d mods to publish to CurseForge", len(cfg.CurseForge.Mods)))
	curseForgePublisher := publisher.NewCurseForgePublisher(*logger, &cfg)
	curseForgePublisher.Publish()
}
