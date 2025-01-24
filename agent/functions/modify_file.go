package functions

import (
	"fmt"
	"os"
	"strings"

	"github.com/clover0/issue-agent/store"
)

const FuncModifyFile = "modify_file"

func InitModifyFileFunction() Function {
	f := Function{
		Name: FuncModifyFile,
		Description: strings.ReplaceAll(`Modify the file at path with the contents of content_text.
 Modified file must be full file content including modified content`, "\n", ""),
		Func: ModifyFile,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path of the file to be modified",
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

	functionsMap[FuncModifyFile] = f

	return f
}

type ModifyFileInput struct {
	Path        string `json:"path"`
	ContentText string `json:"content_text"`
}

func ModifyFile(input ModifyFileInput) (store.File, error) {
	if err := guardPath(input.Path); err != nil {
		return store.File{}, err
	}

	var file store.File
	f, err := os.Create(input.Path)
	if err != nil {
		return file, fmt.Errorf("modify %s: %w", input.Path, err)
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
