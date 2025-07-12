package functions

import (
	"io"
	"os"

	"github.com/clover0/issue-agent/core/store"
)

const FuncOpenFile = "open_file"

func InitOpenFileFunction() Function {
	f := Function{
		Name:        FuncOpenFile,
		Description: "Open the file full content",
		Func:        OpenFile,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The path of the file to open",
				},
			},
			"required":             []string{"path"},
			"additionalProperties": false,
		},
	}

	register(f)

	return f
}

type OpenFileInput struct {
	Path string
}

func OpenFile(input OpenFileInput) (store.File, error) {
	if err := guardPath(input.Path); err != nil {
		return store.File{}, err
	}

	file, err := os.Open(input.Path)
	if err != nil {
		return store.File{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return store.File{}, err
	}

	return store.File{
		Path:    input.Path,
		Content: string(data),
	}, nil
}
