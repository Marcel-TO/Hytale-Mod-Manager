package main

import (
	"fmt"
	"marcel-to/hytale/mod-publisher/config"
	"marcel-to/hytale/mod-publisher/logger"
	"marcel-to/hytale/mod-publisher/publisher"
)

func main() {
	logger := logger.NewLogger("HMP", true)
	logger.Info("Hytale Mod Publisher")

	cfg := config.Config{}
	cfg.GetConfig(*logger, "publish-config.yaml")
	logger.Info("Configuration loaded successfully. Starting to publish....")

	logger.Info("Publishing to CurseForge...")
	logger.Info(fmt.Sprintf("Found %d mods to publish to CurseForge", len(cfg.CurseForge.Mods)))
	curseForgePublisher := publisher.NewCurseForgePublisher(*logger, &cfg)
	for _, mod := range cfg.CurseForge.Mods {
		curseForgePublisher.PublishMod(mod)
	}
}
