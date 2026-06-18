package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"marcel-to/hytale/mod-manager/config"
)

// ReadGradleProperties reads a gradle.properties file and returns mod metadata.
// It delegates parsing to ReadRawModProperties.
func ReadGradleProperties(filePath string, releaseType string) (config.CurseForgeModMetadata, error) {
	modName, version, hytaleVersion, err := ReadRawModProperties(filePath)
	if err != nil {
		return config.CurseForgeModMetadata{}, err
	}
	displayName := fmt.Sprintf("%s-%s+%s", modName, version, hytaleVersion)
	return config.NewCurseForgeModMetadata(displayName, releaseType), nil
}

// UpdateGradleProperties replaces the hytale_version value in a gradle.properties file.
// It writes atomically via a temp file and rename to prevent data loss on crash.
func UpdateGradleProperties(filePath string, newGameVersion string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	r := bufio.NewReader(f)
	var lines []string
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					lines = append(lines, line)
				}
				break
			}
			f.Close()
			return err
		}
		lines = append(lines, line)
	}
	f.Close()

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "hytale_version") {
			if parts := strings.SplitN(trimmed, "=", 2); len(parts) == 2 {
				eqIdx := strings.Index(line, "=")
				lines[i] = line[:eqIdx+1] + newGameVersion + "\n"
			}
			break
		}
	}

	// Write atomically: write to a temp file then rename over the original.
	tmp, err := os.CreateTemp(filepath.Dir(filePath), ".gradle.properties.*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	w := bufio.NewWriter(tmp)
	for _, line := range lines {
		if _, err := w.WriteString(line); err != nil {
			tmp.Close()
			os.Remove(tmpName)
			return err
		}
	}
	if err := w.Flush(); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	return os.Rename(tmpName, filePath)
}
