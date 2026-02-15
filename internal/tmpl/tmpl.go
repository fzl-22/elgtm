package tmpl

import (
	"bytes"
	"fmt"
	"text/template"
)

func Generate[T any](name, content string, data T) (string, error) {
	t, err := template.New(name).Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
