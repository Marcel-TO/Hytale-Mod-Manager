package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/handler"
	"marcel-to/hytale/mod-manager/logger"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type CurseForgePublisher struct {
	BasePublisher
}

func NewCurseForgePublisher(logger logger.Logger, config *config.Config) *CurseForgePublisher {
	return &CurseForgePublisher{
		BasePublisher: BasePublisher{
			Logger: logger,
			url:    "https://legacy.curseforge.com/api/projects/",
			config: config,
		},
	}
}

func (p *CurseForgePublisher) Publish() {
	cache := (&config.UploadCache{}).GetUploadCache(p.Logger)

	for _, mod := range p.config.CurseForge.Mods {
		// Read mod metadata from gradle.properties
		metadata, err := handler.ReadGradleProperties(filepath.Join(mod.RepoLocation, "gradle.properties"), p.Logger)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Failed to read gradle.properties for mod [%s]: %v", mod.Name, err))
			continue
		}

		// Check if the mod version is already cached
		if cache.IsVersionCached(mod.ProjectID, metadata.DisplayName) {
			p.Logger.Info(fmt.Sprintf("Mod [%s] version [%s] already uploaded, skipping", mod.Name, metadata.DisplayName))
			continue
		}

		p.Logger.Info(fmt.Sprintf("Publishing mod [%s] version [%s] to CurseForge...", mod.Name, metadata.DisplayName))
		err = p.PublishMod(mod, metadata)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Failed to publish mod [%s] version [%s]: %v", mod.Name, metadata.DisplayName, err))
		} else {
			cache.AddCacheEntry(mod.ProjectID, metadata.DisplayName)
			cache.SaveUploadCache(p.Logger)
			p.Logger.Info(fmt.Sprintf("Successfully published mod [%s] version [%s] to CurseForge!", mod.Name, metadata.DisplayName))
		}
	}
}

func (p *CurseForgePublisher) PublishMod(mod config.CurseForgeMod, metadata config.CurseForgeModMetadata) error {
	url := p.url + fmt.Sprintf("%d/upload-file", mod.ProjectID)
	method := "POST"

	api_token := os.Getenv("CURSEFORGE_API_TOKEN")
	if api_token == "" {
		return fmt.Errorf("it appears CURSEFORGE_API_TOKEN is not set in .env file")
	}

	// Prepare the multipart form data payload
	payload, contentType, err := preparePayload(metadata, mod.RepoLocation)
	if err != nil {
		return fmt.Errorf("error while preparing the payload: %w", err)
	}

	// Create and send the HTTP request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("error while creating the HTTP request: %w", err)
	}
	req.Header.Add("X-Api-Token", api_token)
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

	// Check if the response body contains the id of the uploaded file
	if !bytes.Contains(body, []byte(`"id":`)) {
		return fmt.Errorf("failed to upload mod: %s", string(body))
	}

	return nil
}

func preparePayload(metadata config.CurseForgeModMetadata, modFilePath string) (*bytes.Buffer, string, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// Add metadata as a form field
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, "", err
	}
	_ = writer.WriteField("metadata", string(metadataJSON))

	// Add the mod file as a form file
	file, err := os.Open(filepath.Join(modFilePath, "build", "libs", fmt.Sprintf("%s.jar", metadata.DisplayName)))
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	part2, err := writer.CreateFormFile("file", filepath.Base(filepath.Join(modFilePath, "build", "libs", fmt.Sprintf("%s.jar", metadata.DisplayName))))
	_, err = io.Copy(part2, file)
	if err != nil {
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return payload, writer.FormDataContentType(), nil
}
