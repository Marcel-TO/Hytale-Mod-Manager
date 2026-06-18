package operations

import (
	"fmt"
	"marcel-to/hytale/mod-manager/builder"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/publisher"
)

func RunPublish(log *logger.Logger, args config.PublishArgumentConfig, cfg *config.Config) error {
	if args.DoBuild {
		log.Info("Building mods before publishing...")
		for _, mod := range cfg.Mods {
			log.Info(fmt.Sprintf("Building mod [%s]...", mod.Name))
			if err := builder.Build(mod.RepoLocation, log); err != nil {
				return fmt.Errorf("failed to build mod [%s]: %w", mod.Name, err)
			}
			log.Info(fmt.Sprintf("Successfully built mod [%s]", mod.Name))
		}
	}

	log.Info("Publishing mods...")
	log.Info(fmt.Sprintf("Found %d mods to publish", len(cfg.Mods)))
	curseForgePublisher := publisher.NewCurseForgePublisher(log)
	curseForgePublisher.Publish(cfg.Mods)
	return nil
}
