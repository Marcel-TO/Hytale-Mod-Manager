package handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"marcel-to/hytale/mod-manager/config"
	"marcel-to/hytale/mod-manager/logger"
)

// ReadGradleProperties reads the content of a Gradle properties file at the specified filePath and returns it as a CurseForgeModMetadata struct.
func ReadGradleProperties(filePath string, logger logger.Logger) (config.CurseForgeModMetadata, error) {
	data, err := os.Open(filePath)
	if err != nil {
		return config.CurseForgeModMetadata{}, err
	}

	defer data.Close()

	r := bufio.NewReader(data)

	mod_name := ""
	version := ""
	hytale_version := ""

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return config.CurseForgeModMetadata{}, err
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lineParts := strings.SplitN(line, "=", 2)
		if len(lineParts) != 2 {
			continue
		}

		key := strings.TrimSpace(lineParts[0])
		value := strings.TrimSpace(lineParts[1])

		switch key {
		case "mod_name":
			mod_name = value
		case "version":
			version = value
		case "hytale_version":
			hytale_version = value
		}

		// If all required metadata is found, stop reading the file
		if mod_name != "" && version != "" && hytale_version != "" {
			break
		}
	}

	if mod_name == "" || version == "" || hytale_version == "" {
		return config.CurseForgeModMetadata{}, fmt.Errorf("missing required metadata in gradle.properties")
	}

	displayName := fmt.Sprintf("%s-%s+%s", mod_name, version, hytale_version)
	return config.NewCurseForgeModMetadata(displayName), nil
}
