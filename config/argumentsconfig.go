package config

import "flag"

type UpdateArgumentConfig struct {
	GameVersion string
	IsPublish   bool
	IsCopying   bool
	Committing  bool
}

type PublishArgumentConfig struct {
	DoBuild bool
}

func ParseUpdateArguments(cmd *flag.FlagSet, args []string) (UpdateArgumentConfig, error) {
	gameVersion := cmd.String("version", "", "The new game version to set in gradle.properties before building and publishing")
	isCopying := cmd.Bool("copy", true, "Indicates whether to copy the built JARs to the Hytale mods directory after building")
	committing := cmd.Bool("commit", true, "Indicates whether to commit the changes to the repository after updating the game version and building")
	isPublish := cmd.Bool("publish", true, "Indicates whether to publish the mods to CurseForge after updating the game version and building")

	if err := cmd.Parse(args); err != nil {
		return UpdateArgumentConfig{}, err
	}

	return UpdateArgumentConfig{
		GameVersion: *gameVersion,
		IsPublish:   *isPublish,
		IsCopying:   *isCopying,
		Committing:  *committing,
	}, nil
}

func ParsePublishArguments(cmd *flag.FlagSet, args []string) (PublishArgumentConfig, error) {
	doBuild := cmd.Bool("build", false, "Build each mod jar before publishing")
	if err := cmd.Parse(args); err != nil {
		return PublishArgumentConfig{}, err
	}

	return PublishArgumentConfig{
		DoBuild: *doBuild,
	}, nil
}
