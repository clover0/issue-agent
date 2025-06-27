package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestEnsureDirAndEnter(t *testing.T) {
	t.Parallel()

	tempBaseDir := t.TempDir()

	// Save the current working directory to restore it after the test
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	// Ensure we return to the original directory after the test
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore original directory: %v", err)
		}
	}()

	// Create a file to test the error case (can't create directory with same name as file)
	filePathForErrorTest := filepath.Join(tempBaseDir, "file_not_dir")
	err = os.WriteFile(filePathForErrorTest, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create file for error test: %v", err)
	}

	tests := map[string]struct {
		dir     string
		wantErr bool
	}{
		"create and enter new directory": {
			dir:     filepath.Join(tempBaseDir, "new_dir"),
			wantErr: false,
		},
		"create and enter nested directory": {
			dir:     filepath.Join(tempBaseDir, "nested", "dir", "path"),
			wantErr: false,
		},
		"enter existing directory": {
			dir:     tempBaseDir,
			wantErr: false,
		},
		"error when directory can't be created": {
			dir:     filepath.Join(filePathForErrorTest, "subdir"), // Can't create dir inside a file
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// We will process sequentially because we need to change the execution directory.
			// t.Parallel()

			// Reset to original directory before each test
			if err := os.Chdir(originalDir); err != nil {
				t.Fatalf("failed to reset to original directory: %v", err)
			}

			err := EnsureDirAndEnter(tt.dir)

			if tt.wantErr {
				assert.HasError(t, err)
				return
			}

			assert.NoError(t, err)

			currentDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get current directory: %v", err)
			}

			assert.Equal(t, filepath.Base(currentDir), filepath.Base(tt.dir))
		})
	}
}
