package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyToHytaleHandler(sourcePath, destinationPath string) error {
	modName, version, hytaleVersion, err := ReadRawModProperties(filepath.Join(sourcePath, "gradle.properties"))
	if err != nil {
		return fmt.Errorf("failed to read gradle.properties: %w", err)
	}

	jarName := fmt.Sprintf("%s-%s+%s.jar", modName, version, hytaleVersion)
	srcJar := filepath.Join(sourcePath, "build", "libs", jarName)

	if _, err := os.Stat(srcJar); err != nil {
		return fmt.Errorf("JAR file not found at %s: %w", srcJar, err)
	}

	entries, err := os.ReadDir(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to read destination directory: %w", err)
	}
	// Remove any previous existing JARs for this mod in the destination directory
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") && strings.Contains(entry.Name(), modName) {
			if err := os.Remove(filepath.Join(destinationPath, entry.Name())); err != nil {
				return fmt.Errorf("failed to remove existing JAR %s: %w", entry.Name(), err)
			}
		}
	}

	return copyJarFile(srcJar, filepath.Join(destinationPath, jarName))
}

func ReadRawModProperties(filePath string) (modName, version, hytaleVersion string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		switch strings.TrimSpace(parts[0]) {
		case "mod_name":
			modName = strings.TrimSpace(parts[1])
		case "version":
			version = strings.TrimSpace(parts[1])
		case "hytale_version":
			hytaleVersion = strings.TrimSpace(parts[1])
		}
		if modName != "" && version != "" && hytaleVersion != "" {
			break
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	if modName == "" || version == "" || hytaleVersion == "" {
		err = fmt.Errorf("missing required metadata in gradle.properties")
	}
	return
}

func copyJarFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
