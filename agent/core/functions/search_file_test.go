package functions_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/clover0/issue-agent/core/functions"
	"github.com/clover0/issue-agent/test/assert"
)

func TestSearchFiles(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	createTestFile := func(path, content string) {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	createSymlink := func(target, linkname string) {
		fullTarget := filepath.Join(tempDir, target)
		fullLinkname := filepath.Join(tempDir, linkname)
		if err := os.Symlink(fullTarget, fullLinkname); err != nil {
			t.Fatalf("failed to create symlink: %v", err)
		}
	}

	createTestFile("file1.txt", "This is a test file containing keyword1.")
	createTestFile("file2.txt", "This is another file with keyword2.")
	createTestFile("subdir/file3.txt", "Subdirectory file also has keyword1.")
	createTestFile(".hidden/file4.txt", "Hidden directory contains keyword1.")
	createTestFile("subdir/.hidden.txt", "This is a hidden file with keyword2.")

	createSymlink("file1.txt", "symlink_to_file1.txt")
	createSymlink("subdir/file3.txt", "symlink_to_file3_in_subdir.txt")

	tests := map[string]struct {
		input    functions.SearchFilesInput
		expected []string
		wantErr  bool
	}{
		"search with keyword1": {
			input: functions.SearchFilesInput{
				Keyword: "keyword1",
				Path:    tempDir,
			},
			expected: []string{
				filepath.Join(tempDir, "file1.txt"),
				filepath.Join(tempDir, "subdir/file3.txt"),
			},
			wantErr: false,
		},
		"search with keyword2": {
			input: functions.SearchFilesInput{
				Keyword: "keyword2",
				Path:    tempDir,
			},
			expected: []string{
				filepath.Join(tempDir, "file2.txt"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
			},
			wantErr: false,
		},
		"search with non-existent keyword": {
			input: functions.SearchFilesInput{
				Keyword: "non-existent-keyword",
				Path:    tempDir,
			},
			expected: make([]string, 0),
			wantErr:  false,
		},
		"search in non-existent path": {
			input: functions.SearchFilesInput{
				Keyword: "keyword1",
				Path:    filepath.Join(tempDir, "not-exist"),
			},
			expected: nil,
			wantErr:  true,
		},
		"search in subdirectory only": {
			input: functions.SearchFilesInput{
				Keyword: "keyword1",
				Path:    filepath.Join(tempDir, "subdir"),
			},
			expected: []string{
				filepath.Join(tempDir, "subdir/file3.txt"),
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := functions.SearchFiles(tt.input)

			if tt.wantErr {
				assert.HasError(t, err)
				return
			}

			slices.Sort(result)
			slices.Sort(tt.expected)

			assert.NoError(t, err)
			assert.Equal(t, len(result), len(tt.expected))
			assert.EqualStringSlices(t, result, tt.expected)
		})
	}
}
