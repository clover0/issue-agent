package common

import (
	"fmt"
	"os"
)

func EnsureDirAndEnter(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("change directory: %w", err)
	}

	return nil
}
