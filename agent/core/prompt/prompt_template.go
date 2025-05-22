package prompt

import (
	"gopkg.in/yaml.v3"
)

type Template struct {
	Agents []AgentPromptTemplate
}

type AgentPromptTemplate struct {
	Name           string `yaml:"name"`
	SystemTemplate string `yaml:"system_prompt"`
	UserTemplate   string `yaml:"user_prompt"`
}

func LoadPrompt() (Template, error) {
	var pt Template

	data := DefaultTemplate()
	err := yaml.Unmarshal(data, &pt)
	if err != nil {
		return pt, err
	}

	return pt, nil
}
