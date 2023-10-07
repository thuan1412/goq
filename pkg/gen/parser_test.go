package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSingleFile(t *testing.T) {
	fpath := "./fixtures/task.go"
	t.Run("parse single file", func(t *testing.T) {
		taskMeta := parseFile(fpath)
		assert.Equal(t, "Greeter", taskMeta.Name)
		assert.Equal(t, "string", taskMeta.PayloadType)
	})
}

func TestParseFileList(t *testing.T) {
	t.Run("parse list of file with pattern", func(t *testing.T) {
		taskMetas := ParseFiles("./**/*task.go")
		assert.Equal(t, 1, len(taskMetas))
		assert.Equal(t, "Greeter", taskMetas[0].Name)
		assert.Equal(t, "string", taskMetas[0].PayloadType)
	})
}
