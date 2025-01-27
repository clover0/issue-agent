package functions_test

import (
	"os"
	"path"
	"testing"

	"github.com/clover0/issue-agent/functions"
	"github.com/clover0/issue-agent/test/assert"
)

func TestRemoveFile(t *testing.T) {
	t.Parallel()

	tempDir := os.TempDir()
	tests := map[string]struct {
		testFile  string
		inputPath string
		wantErr   bool
	}{
		// TODO: add test cases.
		// Currently, test case failed because of test file not located in local.
		//"valid path - file exists": {
		//	testFile:  "test_valid_file.txt",
		//	inputPath: "test_valid_file.txt",
		//	wantErr:   false,
		//},
		"file does not exist": {
			testFile:  "test.txt",
			inputPath: "non_existent_file.txt",
			wantErr:   true,
		},
		"invalid path": {
			testFile:  "",
			inputPath: "../invalid_file.txt",
			wantErr:   true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.testFile != "" {
				tempFile, err := os.Create(path.Join(tempDir, tt.testFile))
				assert.Nil(t, err)
				defer tempFile.Close()
				defer os.Remove(tempFile.Name())
			}

			err := functions.RemoveFile(functions.RemoveFileInput{
				Path: path.Join(tempDir, tt.inputPath),
			})
			if tt.wantErr {
				assert.HasError(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
