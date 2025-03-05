package prompt

import _ "embed"

//go:embed prompt_en.yml
var defaultENPrompt []byte

func DefaultTemplate() []byte {
	return defaultENPrompt
}
