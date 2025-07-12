package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

type PromptTemplate interface {
	Build() (Prompt, error)
}

func ParseTemplate[T any](templateStr string, values T) (string, error) {
	tpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", err
	}

	tplbuff := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(tplbuff, values); err != nil {
		return "", fmt.Errorf(" execute prompt template: %w", err)
	}

	return tplbuff.String(), nil
}
