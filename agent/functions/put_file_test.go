package functions_test

import (
	"os"
	"path/filepath"
	"testing"

	"github/clover0/github-issue-agent/functions"
)

func TestPutFile(t *testing.T) {
	t.Run("successfully writes to a file", func(t *testing.T) {
		path := "testdata/testfile.txt"
		content := "Hello, World!"

		// Clean up before and after the test
		defer os.RemoveAll("testdata")

		input := functions.PutFileInput{
			OutputPath:  path,
			ContentText: content,
		}

		file, err := functions.PutFile(input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if file.Path != path {
			t.Errorf("expected path %s, got %s", path, file.Path)
		}

		if file.Content != content+"\n" {
			t.Errorf("expected content %s, got %s", content+"\n", file.Content)
		}

		// Verify file content
		writtenContent, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if string(writtenContent) != content+"\n" {
			t.Errorf("expected file content %s, got %s", content+"\n", string(writtenContent))
		}
	})

	t.Run("fails with invalid path", func(t *testing.T) {
		input := functions.PutFileInput{
			OutputPath:  "../invalidpath/testfile.txt",
			ContentText: "Invalid path test",
		}

		_, err := functions.PutFile(input)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}
