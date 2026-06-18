package main

import (
	"flag"
	"fmt"
	"marcel-to/hytale/mod-manager/builder"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/operations"
	"marcel-to/hytale/mod-manager/publisher"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.NewLogger("HMP", true)
	log.Info("Hytale Mod Publisher")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.LoadConfig("publish-config.yaml")
	if err != nil {
		log.Error("Failed to load configuration: " + err.Error())
		os.Exit(1)
	}
	log.Info("Configuration loaded successfully!")

	// Load API token from .env file; .env is optional — env vars may be set by the runtime.
	if err := godotenv.Load(); err != nil {
		log.Warning(fmt.Sprintf(".env file not loaded: %v", err))
	}

	switch os.Args[1] {
	case "publish":
		publishCmd := flag.NewFlagSet("publish", flag.ExitOnError)
		doBuild := publishCmd.Bool("build", false, "Build each mod jar before publishing")
		publishCmd.Parse(os.Args[2:])
		if err := runPublish(log, *doBuild, cfg); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		args := config.ParseUpdateArguments(updateCmd, os.Args[2:])
		if err := operations.RunUpdate(log, args, cfg); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runPublish(log *logger.Logger, doBuild bool, cfg *config.Config) error {
	if doBuild {
		log.Info("Building mods before publishing...")
		for _, mod := range cfg.CurseForge.Mods {
			log.Info(fmt.Sprintf("Building mod [%s]...", mod.Name))
			if err := builder.Build(mod.RepoLocation, log); err != nil {
				return fmt.Errorf("failed to build mod [%s]: %w", mod.Name, err)
			}
			log.Info(fmt.Sprintf("Successfully built mod [%s]", mod.Name))
		}
	}

	log.Info("Publishing to CurseForge...")
	log.Info(fmt.Sprintf("Found %d mods to publish to CurseForge", len(cfg.CurseForge.Mods)))
	curseForgePublisher := publisher.NewCurseForgePublisher(log, cfg)
	curseForgePublisher.Publish()
	return nil
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: mod-manager <command> [options]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  publish         Publish all configured mods to CurseForge")
	fmt.Fprintln(os.Stderr, "  update version  Update the game version and publish all configured mods to CurseForge")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Options for 'publish':")
	fmt.Fprintln(os.Stderr, "  --build         Build each mod jar before publishing")
	fmt.Fprintln(os.Stderr, "                  Uses 'just build' if available, otherwise './gradlew build'")
}
