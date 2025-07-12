package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FuncListFiles = "list_files"

func InitListFilesFunction() Function {
	f := Function{
		Name: FuncListFiles,
		Description: strings.ReplaceAll(`List the files within the direc tory like Unix ls command.
Each line contains the file mode, byte size, and name. If you want to list subdirectories recursively, use the depth option.`, "\n", ""),
		Func: ListFiles,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The valid path to list within its directory",
				},
				"depth": map[string]any{
					"type":        "number",
					"description": "The depth of the directory to list subdirectory recursively. Default is 2",
					"minimum":     1,
					"default":     2,
					"maximum":     3,
				},
			},
			"required":             []string{"path"},
			"additionalProperties": false,
		},
	}

	register(f)

	return f
}

type ListFilesInput struct {
	Path  string
	Depth int
}

func ListFiles(input ListFilesInput) ([]string, error) {
	if err := guardPath(input.Path); err != nil {
		return nil, err
	}

	if _, err := os.Stat(input.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist: %w", input.Path, err)
	}

	if input.Depth <= 0 {
		input.Depth = 2
	}

	if input.Depth > 3 {
		input.Depth = 3
	}

	return listFilesRecursive(input.Path, 1, input.Depth, "")
}

func listFilesRecursive(path string, currentDepth, maxDepth int, prefix string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can't read directory at %s: %w", path, err)
	}

	var files []string
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info error: %w", err)
		}

		fullName := entry.Name()
		if prefix != "" {
			fullName = prefix + "/" + entry.Name()
		}

		files = append(files,
			fmt.Sprintf("%s %d %s", info.Mode(), info.Size(), fullName),
		)

		if entry.IsDir() && currentDepth < maxDepth {
			subFiles, err := listFilesRecursive(
				filepath.Join(path, entry.Name()),
				currentDepth+1,
				maxDepth,
				fullName,
			)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}
