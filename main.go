package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
	"marcel-to/hytale/mod-manager/operations"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.NewLogger("HMP", true)
	os.Exit(Run(log, os.Args[1:], os.Stderr))
}

func Run(log *logger.Logger, arguments []string, errorOutput io.Writer) int {
	log.Info("Hytale Mod Publisher")

	if len(arguments) == 0 {
		printUsage(errorOutput)
		return 1
	}

	command := flag.NewFlagSet(arguments[0], flag.ContinueOnError)
	command.SetOutput(errorOutput)
	configPath := command.String("config", "publish-config.yaml", "Path to the YAML configuration file")
	envPath := command.String("env", ".env", "Path to the environment file")

	var runOperation func(*config.Config) error
	var parseErr error
	switch arguments[0] {
	case "publish":
		args, err := config.ParsePublishArguments(command, arguments[1:])
		parseErr = err
		runOperation = func(cfg *config.Config) error {
			return operations.RunPublish(log, args, cfg)
		}
	case "update":
		args, err := config.ParseUpdateArguments(command, arguments[1:])
		parseErr = err
		runOperation = func(cfg *config.Config) error {
			return operations.RunUpdate(log, args, cfg)
		}
	default:
		fmt.Fprintf(errorOutput, "Unknown command: %s\n\n", arguments[0])
		printUsage(errorOutput)
		return 1
	}
	if parseErr != nil {
		if parseErr == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if command.NArg() != 0 {
		fmt.Fprintf(errorOutput, "Unexpected argument: %q\n", command.Arg(0))
		return 2
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Error("Failed to load configuration: " + err.Error())
		return 1
	}
	log.Info("Configuration loaded successfully!")

	if err := godotenv.Load(*envPath); err != nil {
		detail := "invalid environment file syntax"
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			detail = pathErr.Err.Error()
		}
		explicitEnv := false
		command.Visit(func(f *flag.Flag) {
			if f.Name == "env" {
				explicitEnv = true
			}
		})
		if explicitEnv {
			log.Error(fmt.Sprintf("Failed to load environment file %q: %s", *envPath, detail))
			return 1
		}
		log.Warning(fmt.Sprintf(".env file not loaded: %s", detail))
	}

	if err := runOperation(cfg); err != nil {
		log.Error(err.Error())
		return 1
	}
	return 0
}

func printUsage(output io.Writer) {
	fmt.Fprintln(output, "Usage: mod-manager <command> [options]")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Commands:")
	fmt.Fprintln(output, "  publish         Publish all configured mods to CurseForge")
	fmt.Fprintln(output, "  update          Update the game version for all configured mods")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Options for 'publish' and 'update':")
	fmt.Fprintln(output, "  --config <file> Read configuration from this file, default: publish-config.yaml")
	fmt.Fprintln(output, "  --env <file>    Read environment variables from this file, default: .env")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Options for 'publish':")
	fmt.Fprintln(output, "  --build         Build each mod jar before publishing")
	fmt.Fprintln(output, "                  Uses 'just build' if available, otherwise './gradlew build'")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Run 'mod-manager <command> --help' to see all command options.")
}
