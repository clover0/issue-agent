package functions_test

import (
	"os"
	"path"
	"testing"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/core/store"
	"github.com/clover0/issue-agent/test/assert"
)

func TestOpenFile(t *testing.T) {
	t.Parallel()

	tempDir := os.TempDir()
	testContent := "test file content"

	tests := map[string]struct {
		setupFile    string
		setupContent string
		inputPath    string
		wantFile     store.File
		wantErr      bool
	}{
		"valid path - file exists": {
			setupFile:    "test_valid_file.txt",
			setupContent: testContent,
			inputPath:    "test_valid_file.txt",
			wantFile: store.File{
				Path:    path.Join(tempDir, "test_valid_file.txt"),
				Content: testContent,
			},
			wantErr: false,
		},
		"file does not exist": {
			setupFile:    "",
			setupContent: "",
			inputPath:    "non_existent_file.txt",
			wantFile:     store.File{},
			wantErr:      true,
		},
		"invalid path": {
			setupFile:    "",
			setupContent: "",
			inputPath:    "../invalid_file.txt",
			wantFile:     store.File{},
			wantErr:      true,
		},
		"empty file": {
			setupFile:    "empty_file.txt",
			setupContent: "",
			inputPath:    "empty_file.txt",
			wantFile: store.File{
				Path:    path.Join(tempDir, "empty_file.txt"),
				Content: "",
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup test file if needed
			var tempFilePath string
			if tt.setupFile != "" {
				tempFilePath = path.Join(tempDir, tt.setupFile)
				err := os.WriteFile(tempFilePath, []byte(tt.setupContent), 0644)
				assert.Nil(t, err)
				defer func(name string) {
					err := os.Remove(name)
					if err != nil {
						t.Errorf("failed to remove test file %s: %v", name, err)
					}
				}(tempFilePath)
			}

			result, err := functions.OpenFile(functions.OpenFileInput{
				Path: path.Join(tempDir, tt.inputPath),
			})

			if tt.wantErr {
				assert.HasError(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.wantFile.Path, result.Path)
			assert.Equal(t, tt.wantFile.Content, result.Content)
		})
	}
}
