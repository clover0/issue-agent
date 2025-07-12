package functions

import (
	"fmt"
	"os"
)

const FuncRemoveFile = "remove_file"

func InitRemoveFileFunction() Function {
	f := Function{
		Name:        FuncRemoveFile,
		Description: "Remove a file specified by the path",
		Func:        RemoveFile,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path of the file to be removed",
				},
			},
			"required":             []string{"path"},
			"additionalProperties": false,
		},
	}

	register(f)

	return f
}

type RemoveFileInput struct {
	Path string `json:"path"`
}

func RemoveFile(input RemoveFileInput) error {
	if err := guardPath(input.Path); err != nil {
		return err
	}

	if err := os.Remove(input.Path); err != nil {
		return fmt.Errorf("removing %s: %w", input.Path, err)
	}

	return nil
}
