package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

type Template interface {
	Build() (Prompt, error)
}

func ParseTemplate(templateStr string, values any) (string, error) {
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
