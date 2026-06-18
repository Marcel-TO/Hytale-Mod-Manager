package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"marcel-to/hytale/mod-manager/logger"
)

// Constants
const cacheFilePath = "cache/upload_cache.yaml"

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
	ReleaseType  string `yaml:"releaseType"`
}

type CurseForgeModMetadata struct {
	Changelog     string `json:"changelog"`
	ChangelogType string `json:"changelogType"`
	DisplayName   string `json:"displayName"`
	ReleaseType   string `json:"releaseType"`
}

type UploadCache struct {
	CurseForgeUploadCache []CurseForgeUploadCache `yaml:"curseforgeUploadCache"`
}

type CurseForgeUploadCache struct {
	ProjectID int    `yaml:"projectId"`
	Version   string `yaml:"version"`
}

func NewCurseForgeModMetadata(displayName string) CurseForgeModMetadata {
	return CurseForgeModMetadata{
		DisplayName:   displayName,
		ReleaseType:   "release",
		Changelog:     "",
		ChangelogType: "markdown",
	}
}

func (cache *UploadCache) GetUploadCache(logger logger.Logger) *UploadCache {
	yamlFile, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cache
		}
		logger.Error("Failed to read upload cache: " + err.Error())
		return cache
	}

	err = yaml.Unmarshal(yamlFile, cache)
	if err != nil {
		logger.Error("Failed to parse upload cache: " + err.Error())
	}
	return cache
}

func (cache *UploadCache) SaveUploadCache(logger logger.Logger) {
	data, err := yaml.Marshal(cache)
	if err != nil {
		logger.Error("Failed to serialize upload cache: " + err.Error())
		return
	}
	err = os.WriteFile(cacheFilePath, data, 0644)
	if err != nil {
		logger.Error("Failed to write upload cache: " + err.Error())
	}
}

func (cache *UploadCache) IsVersionCached(projectID int, version string) bool {
	for _, entry := range cache.CurseForgeUploadCache {
		if entry.ProjectID == projectID && entry.Version == version {
			return true
		}
	}
	return false
}

func (cache *UploadCache) AddCacheEntry(projectID int, version string) {
	if cache.IsVersionCached(projectID, version) {
		return
	}

	// Overwrite existing entry for the same project ID if it exists
	for i, entry := range cache.CurseForgeUploadCache {
		if entry.ProjectID == projectID {
			cache.CurseForgeUploadCache = append(cache.CurseForgeUploadCache[:i], cache.CurseForgeUploadCache[i+1:]...)
			break
		}
	}

	cache.CurseForgeUploadCache = append(cache.CurseForgeUploadCache, CurseForgeUploadCache{
		ProjectID: projectID,
		Version:   version,
	})
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
