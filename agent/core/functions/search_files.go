package functions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const FuncSearchFiles = "search_files"

func InitSearchFilesFunction() Function {
	f := Function{
		Name: FuncSearchFiles,
		Description: strings.ReplaceAll(`Search for files containing specific keyword (e.g., "xxx")
 within a directory path recursively`, "\n", ""),
		Func: SearchFiles,
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"keyword": map[string]interface{}{
					"type":        "string",
					"description": "The keyword to search for.",
				},
				"path": map[string]interface{}{
					"type":        "string",
					"description": "The path to search within its directory",
				},
			},
			"required":             []string{"keyword", "path"},
			"additionalProperties": false,
		},
	}

	functionsMap[FuncSearchFiles] = f

	return f
}

type SearchFilesInput struct {
	Keyword string
	Path    string
}

func SearchFiles(input SearchFilesInput) ([]string, error) {
	if err := guardPath(input.Path); err != nil {
		return nil, err
	}

	if _, err := os.Stat(input.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist: %w", input.Path, err)
	}

	var currentDirs = []string{".", "./"}
	var fileNames []string
	err := filepath.WalkDir(input.Path, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			// Skip hidden directories
			if !slices.Contains(currentDirs, d.Name()) && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		f, fileErr := os.Open(path)
		if fileErr != nil {
			return fmt.Errorf("failed to open file: %w", fileErr)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, input.Keyword) {
				fileNames = append(fileNames, filepath.Clean(path))
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return fileNames, nil
}
