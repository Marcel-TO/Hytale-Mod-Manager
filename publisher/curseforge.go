package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"marcel-to/hytale/mod-publisher/config"
	"marcel-to/hytale/mod-publisher/handler"
	"marcel-to/hytale/mod-publisher/logger"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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

func (p *CurseForgePublisher) PublishMod(mod config.CurseForgeMod) {
	url := p.url + fmt.Sprintf("%d/upload-file", mod.ProjectID)
	method := "POST"

	metadata, err := handler.ReadGradleProperties(filepath.Join(mod.RepoLocation, "gradle.properties"), p.Logger)
	if err != nil {
		fmt.Println(err)
		return
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// Add metadata as a form field
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = writer.WriteField("metadata", string(metadataJSON))

	// Add the mod file as a form file
	file, err := os.Open(filepath.Join(mod.RepoLocation, "build", "libs", fmt.Sprintf("%s.jar", metadata.DisplayName)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	part2, err := writer.CreateFormFile("file", filepath.Base(filepath.Join(mod.RepoLocation, "build", "libs", fmt.Sprintf("%s.jar", metadata.DisplayName))))
	_, err = io.Copy(part2, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = godotenv.Load()
	if err != nil {
		p.Logger.Error("Error loading .env file")
		return
	}

	api_token := os.Getenv("CURSEFORGE_API_TOKEN")
	if api_token == "" {
		p.Logger.Error("CURSEFORGE_API_TOKEN not set in .env file")
		return
	}

	req.Header.Add("X-Api-Token", api_token)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
