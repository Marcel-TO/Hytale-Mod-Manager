package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"marcel-to/hytale/mod-manager/handler"
	"marcel-to/hytale/mod-manager/logger"
)

// Build attempts to build the mod jar in the given repoLocation.
// It prefers `just build` if the `just` binary is available AND a justfile
// exists in repoLocation, otherwise falls back to `./gradlew build`
func Build(repoLocation string, log logger.Logger) error {
	if hasJust() && hasJustfile(repoLocation) {
		log.Info(fmt.Sprintf("Running 'just build' in %s", repoLocation))
		return run(repoLocation, "just", "build")
	}

	gradlew := "./gradlew"
	log.Info(fmt.Sprintf("'just' not found or no justfile present, running '%s build' in %s", gradlew, repoLocation))
	return run(repoLocation, gradlew, "build")
}

func UpdateGameVersion(repoLocation string, gameVersion string, log logger.Logger) error {
	err := handler.UpdateGradleProperties(filepath.Join(repoLocation, "gradle.properties"), gameVersion, log)
	if err != nil {
		return fmt.Errorf("failed to update gradle.properties: %w", err)
	}

	return Build(repoLocation, log)
}

// hasJust reports whether the 'just' binary is available in the system's PATH.
func hasJust() bool {
	_, err := exec.LookPath("just")
	return err == nil
}

// hasJustfile reports whether repoLocation contains a justfile.
// just recognises both 'justfile' and 'Justfile' as valid names.
func hasJustfile(repoLocation string) bool {
	for _, name := range []string{"justfile", "Justfile"} {
		if _, err := os.Stat(filepath.Join(repoLocation, name)); err == nil {
			return true
		}
	}
	return false
}

// run executes the given command in the specified directory, streaming its output to stdout and stderr.
func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command '%s %s' failed: %w", name, strings.Join(args, " "), err)
	}
	return nil
}
