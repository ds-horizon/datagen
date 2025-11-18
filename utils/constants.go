package utils

import (
	"fmt"

	"github.com/elliotchance/orderedmap/v3"
)

const (
	DgDirDelimeter       = "___DGDIRDELIM___"
	DefaultMetadataCount = 1
	EncodedBinaryName    = "datagen"
	CompilerBinaryName   = "datagenc"
)

type DgDir struct {
	Name     string
	Models   *orderedmap.OrderedMap[string, []byte]
	Children []*DgDir
}

func (n *DgDir) ModelCount() int {
	total := n.Models.Len()
	for _, child := range n.Children {
		total += child.ModelCount()
	}
	return total
}

func (n *DgDir) prettyPrint(indent string) {
	fmt.Println(indent + n.Name)

	if n.Models.Len() > 0 {
		for k := range n.Models.AllFromFront() {
			fmt.Println(indent + "  [model] " + k)
		}
	}

	for _, child := range n.Children {
		child.prettyPrint(indent + "  ")
	}
}

func (n *DgDir) PrettyPrint() {
	n.prettyPrint("")
}
