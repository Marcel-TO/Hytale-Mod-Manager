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
	Mods []ModConfig `yaml:"mods"`
}

// ModConfig holds mod identity and build information, shared across all platforms.
type ModConfig struct {
	Name         string          `yaml:"name"`
	RepoLocation string          `yaml:"repoLocation"`
	ReleaseType  string          `yaml:"releaseType"`
	Platforms    PlatformTargets `yaml:"platforms"`
}

// PlatformTargets holds optional publishing configuration for each supported platform.
// A nil pointer means the mod should not be published to that platform.
type PlatformTargets struct {
	CurseForge *CurseForgeTarget `yaml:"curseforge,omitempty"`
}

// CurseForgeTarget holds CurseForge-specific publishing configuration.
type CurseForgeTarget struct {
	ProjectID int `yaml:"projectId"`
}

type CurseForgeModMetadata struct {
	Changelog     string `json:"changelog"`
	ChangelogType string `json:"changelogType"`
	DisplayName   string `json:"displayName"`
	ReleaseType   string `json:"releaseType"`
}

// UploadCache tracks which mod versions have already been uploaded to each platform.
type UploadCache struct {
	CurseForge []CurseForgeUploadCacheEntry `yaml:"curseforge,omitempty"`
}

type CurseForgeUploadCacheEntry struct {
	ProjectID int    `yaml:"projectId"`
	Version   string `yaml:"version"`
}

func NewCurseForgeModMetadata(displayName, releaseType string) CurseForgeModMetadata {
	return CurseForgeModMetadata{
		DisplayName:   displayName,
		ReleaseType:   releaseType,
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
	for _, entry := range cache.CurseForge {
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

	for i, entry := range cache.CurseForge {
		if entry.ProjectID == projectID {
			cache.CurseForge = append(cache.CurseForge[:i], cache.CurseForge[i+1:]...)
			break
		}
	}

	cache.CurseForge = append(cache.CurseForge, CurseForgeUploadCacheEntry{
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
