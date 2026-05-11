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

func UpdateGradleProperties(filePath string, newGameVersion string, logger logger.Logger) error {
	data, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer data.Close()

	r := bufio.NewReader(data)
	var lines []string

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return err
			}
			if line != "" {
				lines = append(lines, line)
			}
			break
		}
		lines = append(lines, line)
	}

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "hytale_version") {
			lineParts := strings.SplitN(trimmedLine, "=", 2)
			if len(lineParts) == 2 {
				prefix := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = prefix + "hytale_version=" + newGameVersion + "\n"
			}
			break
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := w.WriteString(line)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
