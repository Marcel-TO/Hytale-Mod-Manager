package git

import (
	"fmt"
	"path/filepath"

	"marcel-to/hytale/mod-manager/builder"
	"marcel-to/hytale/mod-manager/handler"
)

func CommitToGitHandler(repoPath string) error {
	modName, version, hytaleVersion, err := handler.ReadRawModProperties(filepath.Join(repoPath, "gradle.properties"))
	if err != nil {
		return fmt.Errorf("failed to read gradle.properties: %w", err)
	}

	commitMessage := fmt.Sprintf("feat: update %s to version %s+%s", modName, version, hytaleVersion)
	err = builder.Run(repoPath, "git", "add", ".")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}
	err = builder.Run(repoPath, "git", "commit", "-m", commitMessage)
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	err = builder.Run(repoPath, "git", "push")
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	return nil
}
