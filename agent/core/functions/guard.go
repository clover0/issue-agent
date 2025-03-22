package functions

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func guardPath(path string) error {
	if testing.Testing() {
		return nil
	}

	return guardPathInner(path)
}

func guardPathInner(path string) error {
	cleanPath := filepath.Clean(path)
	if !filepath.IsLocal(cleanPath) {
		return fmt.Errorf("path %s is not a local path", path)
	}

	if strings.HasPrefix(cleanPath, "..") {
		return fmt.Errorf("path %s attempts to access parent directory", path)
	}

	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path %s contains not allowed '..'", path)
	}

	if strings.Contains(cleanPath, "~") {
		return fmt.Errorf("path %s contains not allowed '~'", path)
	}

	if strings.Contains(cleanPath, "//") {
		return fmt.Errorf("path %s contains not allowed '//'", path)
	}

	if strings.HasPrefix(cleanPath, "/") {
		return fmt.Errorf("path %s starts with '/', not allowed", path)
	}

	return nil
}
