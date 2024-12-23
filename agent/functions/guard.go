package functions

import (
	"fmt"
	"path/filepath"
	"strings"
)

// guardPath checks if the given path is safe to access within the current directory context
// It prevents path traversal attacks and ensures the path is within allowed boundaries
func guardPath(path string) error {
	// Basic validation checks
	if strings.Contains(path, "..") {
		return fmt.Errorf("path %s contains not allowed '..'", path)
	}
	if strings.Contains(path, "~") {
		return fmt.Errorf("path %s contains not allowed '~'", path)
	}
	if strings.Contains(path, "//") {
		return fmt.Errorf("path %s contains not allowed '//'", path)
	}
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("path %s starts with '/', not allowed", path)
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(path)

	// Ensure the path is local (doesn't try to escape current directory)
	if !filepath.IsLocal(cleanPath) {
		return fmt.Errorf("path %s is not a local path", path)
	}

	// Additional check to ensure the cleaned path doesn't start with ".."
	if strings.HasPrefix(cleanPath, "..") {
		return fmt.Errorf("path %s attempts to access parent directory", path)
	}

	return nil
}