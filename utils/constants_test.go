package utils

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/stretchr/testify/assert"
)

func TestDgDir_ModelCount(t *testing.T) {
	dgDir := &DgDir{
		Name:     "root",
		Models:   orderedmap.NewOrderedMap[string, []byte](),
		Children: []*DgDir{},
	}
	assert.Equal(t, 0, dgDir.ModelCount())

	dgDir.Models.Set("model1", []byte("content"))
	assert.Equal(t, 1, dgDir.ModelCount())

	childModels := orderedmap.NewOrderedMap[string, []byte]()
	childModels.Set("child_model", []byte("content"))
	dgDir.Children = []*DgDir{
		{
			Name:     "child",
			Models:   childModels,
			Children: []*DgDir{},
		},
	}
	assert.Equal(t, 2, dgDir.ModelCount())
}

func TestDgDir_PrettyPrint(t *testing.T) {
	captureOutput := func(fn func()) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		fn()

		_ = w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		return buf.String()
	}

	t.Run("empty directory", func(t *testing.T) {
		dgDir := &DgDir{
			Name:     "root",
			Models:   orderedmap.NewOrderedMap[string, []byte](),
			Children: []*DgDir{},
		}
		output := captureOutput(func() {
			dgDir.PrettyPrint()
		})
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Equal(t, 1, len(lines))
		assert.Equal(t, "root", lines[0])
	})

	t.Run("with models", func(t *testing.T) {
		dgDir := &DgDir{
			Name:     "root",
			Models:   orderedmap.NewOrderedMap[string, []byte](),
			Children: []*DgDir{},
		}
		dgDir.Models.Set("model1", []byte("content"))
		output := captureOutput(func() {
			dgDir.PrettyPrint()
		})
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.GreaterOrEqual(t, len(lines), 2)
		assert.Equal(t, "root", lines[0])
		assert.Contains(t, output, "[model] model1")
	})

	t.Run("with nested children", func(t *testing.T) {
		dgDir := &DgDir{
			Name:     "root",
			Models:   orderedmap.NewOrderedMap[string, []byte](),
			Children: []*DgDir{},
		}
		dgDir.Models.Set("model1", []byte("content"))
		childModels := orderedmap.NewOrderedMap[string, []byte]()
		childModels.Set("child_model", []byte("content"))
		dgDir.Children = []*DgDir{
			{
				Name:     "child",
				Models:   childModels,
				Children: []*DgDir{},
			},
		}
		output := captureOutput(func() {
			dgDir.PrettyPrint()
		})
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.GreaterOrEqual(t, len(lines), 3)
		assert.Equal(t, "root", lines[0])
		assert.Contains(t, output, "[model] model1")
		assert.Contains(t, output, "child")
		assert.Contains(t, output, "[model] child_model")
	})
}
