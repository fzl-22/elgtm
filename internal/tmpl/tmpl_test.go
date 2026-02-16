package tmpl_test

import (
	"testing"

	"github.com/fzl-22/elgtm/internal/tmpl"
	"github.com/stretchr/testify/assert"
)

func TestGenerate_Success(t *testing.T) {
	data := struct{ Name string }{Name: "ELGTM"}
	content := "Hello {{.Name}}!"

	res, err := tmpl.Generate("test", content, data)

	assert.NoError(t, err)
	assert.Equal(t, "Hello ELGTM!", res)
}

func TestGenerate_ParseError(t *testing.T) {
	data := struct{ Name string }{Name: "ELGTM"}
	content := "Hello {{.Name}"

	_, err := tmpl.Generate("test", content, data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse template")
}

func TestGenerate_ExecuteError(t *testing.T) {
	data := struct{ Name string }{Name: "ELGTM"}
	content := "Hello {{.NonExistentValue}}"

	_, err := tmpl.Generate("test", content, data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute template")
}
