package tmpl_test

import (
	"testing"

	"github.com/fzl-22/elgtm/internal/tmpl"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("Success_ValidTemplate", func(t *testing.T) {
		data := struct{ Name string }{Name: "ELGTM"}
		content := "Hello {{.Name}}!"

		res, err := tmpl.Generate("test", content, data)

		assert.NoError(t, err)
		assert.Equal(t, "Hello ELGTM!", res)
	})

	t.Run("Failure_InvalidTemplateSyntax", func(t *testing.T) {
		data := struct{ Name string }{Name: "ELGTM"}
		content := "Hello {{.Name}"

		_, err := tmpl.Generate("test", content, data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse template")
	})

	t.Run("Failure_InvalidTemplateData", func(t *testing.T) {
		data := struct{ Name string }{Name: "ELGTM"}
		content := "Hello {{.NonExistentValue}}"

		_, err := tmpl.Generate("test", content, data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute template")
	})
}
