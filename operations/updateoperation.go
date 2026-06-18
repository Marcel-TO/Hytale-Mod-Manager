package operations

import (
	"fmt"
	"os"
	"path/filepath"

	"marcel-to/hytale/mod-manager/builder"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/git"
	"marcel-to/hytale/mod-manager/handler"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/publisher"
)

func RunUpdate(log *logger.Logger, args config.UpdateArgumentConfig, cfg *config.Config) error {
	for _, mod := range cfg.CurseForge.Mods {
		log.Info(fmt.Sprintf("Processing mod [%s]...", mod.Name))

		// Update the game version and build the mods
		if err := updateGameVersion(log, mod, args.GameVersion); err != nil {
			return err
		}

		// Copy the built JARs to the Hytale mods directory
		if args.IsCopying {
			if err := copyBuiltJarsToHytale(log, mod); err != nil {
				return err
			}
		}

		// Commit the changes to git
		if args.Committing {
			if err := commitChangesToGit(log, mod); err != nil {
				return err
			}
		}
	}

	// Publish the mods to CurseForge after updating the game version and building
	if args.IsPublish {
		log.Info("Publishing to CurseForge...")
		log.Info(fmt.Sprintf("Found %d mods to publish to CurseForge", len(cfg.CurseForge.Mods)))
		curseForgePublisher := publisher.NewCurseForgePublisher(*log, cfg)
		curseForgePublisher.Publish()
	}
	return nil
}

func updateGameVersion(log *logger.Logger, mod config.CurseForgeMod, newVersion string) error {
	// Update the game version and build the mod
	log.Info(fmt.Sprintf("Updating game version to %s and building mod [%s]...", newVersion, mod.Name))
	if err := builder.UpdateGameVersion(mod.RepoLocation, newVersion, *log); err != nil {
		return fmt.Errorf("failed to update and build mod [%s]: %w", mod.Name, err)
	}
	log.Info(fmt.Sprintf("Successfully updated and built mod [%s]", mod.Name))

	return nil
}

func copyBuiltJarsToHytale(log *logger.Logger, mod config.CurseForgeMod) error {
	log.Info(fmt.Sprintf("Copying built JARs for mod [%s] to Hytale mods directory...", mod.Name))

	hytalePath := os.Getenv("HYTALE_PATH")
	if hytalePath == "" {
		return fmt.Errorf("HYTALE_PATH environment variable is not set. Please set it to the root directory of your Hytale installation")
	}
	if err := handler.CopyToHytaleHandler(mod.RepoLocation, filepath.Join(hytalePath, "UserData/Mods"), mod.ReleaseType); err != nil {
		return fmt.Errorf("failed to copy mod [%s] to Hytale mods directory: %w", mod.Name, err)
	}
	log.Info(fmt.Sprintf("Successfully copied mod [%s] to Hytale mods directory", mod.Name))

	return nil
}

func commitChangesToGit(log *logger.Logger, mod config.CurseForgeMod) error {
	if err := git.CommitToGitHandler(mod.RepoLocation); err != nil {
		return fmt.Errorf("failed to commit changes for mod [%s]: %w", mod.Name, err)
	}

	log.Info(fmt.Sprintf("Successfully committed the changes for mod [%s] to git.", mod.Name))
	return nil
}
