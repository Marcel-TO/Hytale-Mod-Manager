package main

import (
	"flag"
	"fmt"
	"marcel-to/hytale/mod-manager/builder"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/publisher"
	"os"
)

func main() {
	log := logger.NewLogger("HMP", true)
	log.Info("Hytale Mod Publisher")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.Config{}
	cfg.GetConfig(*log, "publish-config.yaml")
	log.Info("Configuration loaded successfully!")

	switch os.Args[1] {
	case "publish":
		publishCmd := flag.NewFlagSet("publish", flag.ExitOnError)
		doBuild := publishCmd.Bool("build", false, "Build each mod jar before publishing")
		publishCmd.Parse(os.Args[2:])
		runPublish(log, *doBuild, &cfg)
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		gameVersion := updateCmd.String("version", "", "The new game version to set in gradle.properties before building and publishing")
		updateCmd.Parse(os.Args[2:])
		runUpdate(log, *gameVersion, &cfg)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runPublish(log *logger.Logger, doBuild bool, cfg *config.Config) {
	if doBuild {
		log.Info("Building mods before publishing...")
		for _, mod := range cfg.CurseForge.Mods {
			log.Info(fmt.Sprintf("Building mod [%s]...", mod.Name))
			if err := builder.Build(mod.RepoLocation, *log); err != nil {
				log.Error(fmt.Sprintf("Failed to build mod [%s]: %v", mod.Name, err))
				os.Exit(1)
			}
			log.Info(fmt.Sprintf("Successfully built mod [%s]", mod.Name))
		}
	}

	log.Info("Publishing to CurseForge...")
	log.Info(fmt.Sprintf("Found %d mods to publish to CurseForge", len(cfg.CurseForge.Mods)))
	curseForgePublisher := publisher.NewCurseForgePublisher(*log, cfg)
	curseForgePublisher.Publish()
}

func runUpdate(log *logger.Logger, gameVersion string, cfg *config.Config) {
	log.Info(fmt.Sprintf("Updating game version to %s and building mods...", gameVersion))
	for _, mod := range cfg.CurseForge.Mods {
		log.Info(fmt.Sprintf("Updating and building mod [%s]...", mod.Name))
		if err := builder.UpdateGameVersion(mod.RepoLocation, gameVersion, *log); err != nil {
			log.Error(fmt.Sprintf("Failed to update and build mod [%s]: %v", mod.Name, err))
			os.Exit(1)
		}
		log.Info(fmt.Sprintf("Successfully updated and built mod [%s]", mod.Name))
	}

	log.Info("Publishing to CurseForge...")
	log.Info(fmt.Sprintf("Found %d mods to publish to CurseForge", len(cfg.CurseForge.Mods)))
	curseForgePublisher := publisher.NewCurseForgePublisher(*log, cfg)
	curseForgePublisher.Publish()
}

func printUsage() {
	fmt.Println("Usage: mod-manager <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  publish         Publish all configured mods to CurseForge")
	fmt.Println("  update version  Update the game version and publish all configured mods to CurseForge")
	fmt.Println()
	fmt.Println("Options for 'publish':")
	fmt.Println("  --build         Build each mod jar before publishing")
	fmt.Println("                  Uses 'just build' if available, otherwise './gradlew build'")
}
