package operations

import (
	"fmt"
	"strings"

	"marcel-to/hytale/mod-manager/builder"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/publisher"
)

func RunPublish(log *logger.Logger, args config.PublishArgumentConfig, cfg *config.Config) error {
	modsToPublish := cfg.Mods
	var failedMods []string

	if args.DoBuild {
		log.Info("Building mods before publishing...")
		modsToPublish = nil
		for _, mod := range cfg.Mods {
			log.Info(fmt.Sprintf("Building mod [%s]...", mod.Name))
			if err := builder.Build(mod.RepoLocation, log); err != nil {
				log.Error(fmt.Sprintf("Failed to build mod [%s]: %s", mod.Name, err.Error()))
				log.Warning(fmt.Sprintf("Skipping mod [%s] due to build error. Continuing with next mod...", mod.Name))
				failedMods = append(failedMods, mod.Name)
				continue
			}
			log.Info(fmt.Sprintf("Successfully built mod [%s]", mod.Name))
			modsToPublish = append(modsToPublish, mod)
		}
	}

	if len(modsToPublish) > 0 {
		log.Info("Publishing mods...")
		log.Info(fmt.Sprintf("Found %d mods to publish", len(modsToPublish)))
		curseForgePublisher := publisher.NewCurseForgePublisher(log)
		curseForgePublisher.Publish(modsToPublish)
	}

	if len(failedMods) > 0 {
		return fmt.Errorf("the following mods failed to build and were skipped: %s", strings.Join(failedMods, ", "))
	}

	return nil
}
