package functions_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/clover0/github-issue-agent/functions"
)

func TestPutFile(t *testing.T) {
	t.Run("successfully writes content to file", func(t *testing.T) {
		outputPath := "testdata/testfile.txt"
		content := "Hello, World!"

		defer os.Remove(outputPath)

		file, err := functions.PutFile(functions.PutFileInput{
			OutputPath:  outputPath,
			ContentText: content,
		})

		require.NoError(t, err)
		assert.Equal(t, outputPath, file.Path)
		assert.Equal(t, content+"\n", file.Content)

		actualContent, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		assert.Equal(t, content+"\n", string(actualContent))
	})

	t.Run("fails with invalid path", func(t *testing.T) {
		outputPath := "invalid_path/testfile.txt"
		content := "Hello, World!"

		_, err := functions.PutFile(functions.PutFileInput{
			OutputPath:  outputPath,
			ContentText: content,
		})

		assert.Error(t, err)
	})
}
