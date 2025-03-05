package prompt

import (
	"gopkg.in/yaml.v3"
)

type PromptTemplate struct {
	Agents []struct {
		Name           string `yaml:"name"`
		SystemTemplate string `yaml:"system_prompt"`
		UserTemplate   string `yaml:"user_prompt"`
	}
}

func LoadPrompt() (PromptTemplate, error) {
	var pt PromptTemplate

	data := DefaultTemplate()
	err := yaml.Unmarshal(data, &pt)
	if err != nil {
		return pt, err
	}

	return pt, nil
}
