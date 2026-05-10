package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"marcel-to/hytale/mod-publisher/logger"
)

type Config struct {
	CurseForge CurseForgeConfig `yaml:"curseforge"`
}

type CurseForgeConfig struct {
	Mods []CurseForgeMod `yaml:"mods"`
}

type CurseForgeMod struct {
	Name         string `yaml:"name"`
	ProjectID    int    `yaml:"projectId"`
	RepoLocation string `yaml:"repoLocation"`
}

type CurseForgeModMetadata struct {
	Changelog     string `json:"changelog"`
	ChangelogType string `json:"changelogType"`
	DisplayName   string `json:"displayName"`
	ReleaseType   string `json:"releaseType"`
}

func NewCurseForgeModMetadata(displayName string) CurseForgeModMetadata {
	return CurseForgeModMetadata{
		DisplayName:   displayName,
		ReleaseType:   "release",
		Changelog:     "",
		ChangelogType: "markdown",
	}
}

func (config *Config) GetConfig(logger logger.Logger, filePath string) *Config {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Failed to read config.yaml: " + err.Error())
		return nil
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		logger.Error("Failed to parse config.yaml: " + err.Error())
		return nil
	}

	return config
}
