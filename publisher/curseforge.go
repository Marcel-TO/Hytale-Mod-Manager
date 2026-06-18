package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/handler"
	"marcel-to/hytale/mod-manager/logger"
)

type CurseForgePublisher struct {
	BasePublisher
}

func NewCurseForgePublisher(log *logger.Logger) *CurseForgePublisher {
	return &CurseForgePublisher{
		BasePublisher: BasePublisher{
			Logger: log,
			url:    "https://legacy.curseforge.com/api/projects/",
		},
	}
}

func (p *CurseForgePublisher) Publish(mods []config.ModConfig) {
	cache := (&config.UploadCache{}).GetUploadCache(p.Logger)

	for _, mod := range mods {
		if mod.Platforms.CurseForge == nil {
			continue
		}
		target := mod.Platforms.CurseForge

		metadata, err := handler.ReadGradleProperties(filepath.Join(mod.RepoLocation, "gradle.properties"), mod.ReleaseType)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Failed to read gradle.properties for mod [%s]: %v", mod.Name, err))
			continue
		}

		if cache.IsVersionCached(target.ProjectID, metadata.DisplayName) {
			p.Logger.Info(fmt.Sprintf("Mod [%s] version [%s] already uploaded, skipping", mod.Name, metadata.DisplayName))
			continue
		}

		p.Logger.Info(fmt.Sprintf("Publishing mod [%s] version [%s] to CurseForge...", mod.Name, metadata.DisplayName))
		err = p.PublishMod(target.ProjectID, mod.RepoLocation, metadata)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Failed to publish mod [%s] version [%s]: %v", mod.Name, metadata.DisplayName, err))
		} else {
			cache.AddCacheEntry(target.ProjectID, metadata.DisplayName)
			cache.SaveUploadCache(p.Logger)
			p.Logger.Info(fmt.Sprintf("Successfully published mod [%s] version [%s] to CurseForge!", mod.Name, metadata.DisplayName))
		}
	}
}

func (p *CurseForgePublisher) PublishMod(projectID int, repoLocation string, metadata config.CurseForgeModMetadata) error {
	url := p.url + fmt.Sprintf("%d/upload-file", projectID)

	apiToken := os.Getenv("CURSEFORGE_API_TOKEN")
	if apiToken == "" {
		return fmt.Errorf("CURSEFORGE_API_TOKEN is not set")
	}

	payload, contentType, err := preparePayload(metadata, repoLocation)
	if err != nil {
		return fmt.Errorf("error while preparing the payload: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return fmt.Errorf("error while creating the HTTP request: %w", err)
	}
	req.Header.Add("X-Api-Token", apiToken)
	req.Header.Set("Content-Type", contentType)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending the HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error while reading the response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("CurseForge API returned HTTP %d: %s", res.StatusCode, string(body))
	}

	if !bytes.Contains(body, []byte(`"id":`)) {
		return fmt.Errorf("unexpected response from CurseForge API: %s", string(body))
	}

	return nil
}

func preparePayload(metadata config.CurseForgeModMetadata, modFilePath string) (*bytes.Buffer, string, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, "", err
	}
	if err := writer.WriteField("metadata", string(metadataJSON)); err != nil {
		return nil, "", err
	}

	jarPath := filepath.Join(modFilePath, "build", "libs", fmt.Sprintf("%s.jar", metadata.DisplayName))
	file, err := os.Open(jarPath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(jarPath))
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, "", err
	}
	if err = writer.Close(); err != nil {
		return nil, "", err
	}

	return payload, writer.FormDataContentType(), nil
}
