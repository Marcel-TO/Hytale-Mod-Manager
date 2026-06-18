package config

import (
	"fmt"
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

func (cache *UploadCache) GetUploadCache(log *logger.Logger) *UploadCache {
	yamlFile, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cache
		}
		log.Error("Failed to read upload cache: " + err.Error())
		return cache
	}

	err = yaml.Unmarshal(yamlFile, cache)
	if err != nil {
		log.Error("Failed to parse upload cache: " + err.Error())
	}
	return cache
}

func (cache *UploadCache) SaveUploadCache(log *logger.Logger) {
	data, err := yaml.Marshal(cache)
	if err != nil {
		log.Error("Failed to serialize upload cache: " + err.Error())
		return
	}
	err = os.WriteFile(cacheFilePath, data, 0644)
	if err != nil {
		log.Error("Failed to write upload cache: " + err.Error())
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

// LoadConfig reads and parses the YAML configuration file at filePath.
func LoadConfig(filePath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var cfg Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	return &cfg, nil
}
