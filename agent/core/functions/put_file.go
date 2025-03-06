package functions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clover0/issue-agent/core/store"
)

const FuncPutFile = "put_file"

func InitPutFileFunction() Function {
	f := Function{
		Name:        FuncPutFile,
		Description: "Put new file content to path",
		Func:        PutFile,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path of the file to be changed to the new content",
				},
				"content_text": map[string]interface{}{
					"type":        "string",
					"description": "The new content of the file",
				},
			},
			"required":             []string{"path", "content_text"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncPutFile] = f

	return f
}

type PutFileInput struct {
	Path        string `json:"path"`
	ContentText string `json:"content_text"`
}

func PutFile(input PutFileInput) (store.File, error) {
	if err := guardPath(input.Path); err != nil {
		return store.File{}, err
	}

	var file store.File
	baseDir := filepath.Dir(input.Path)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return file, fmt.Errorf("mkdir all %s error: %w", baseDir, err)
	}

	f, err := os.Create(input.Path)
	if err != nil {
		return file, fmt.Errorf("putting %s: %w", input.Path, err)
	}
	defer f.Close()

	// EOF should be a newline
	if len(input.ContentText) != 0 && input.ContentText[len(input.ContentText)-1] != '\n' {
		input.ContentText += "\n"
	}

	if _, err := f.WriteString(input.ContentText); err != nil {
		return file, err
	}

	return store.File{
		Path:    input.Path,
		Content: input.ContentText,
	}, nil
}
