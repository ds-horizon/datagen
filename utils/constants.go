package utils

import "fmt"

const DgDirDelimeter = "___DGDIRDELIM___"
const DefaultMetadataCount = 1

type DgDir struct {
	Name     string
	Models   map[string][]byte
	Children []*DgDir
}

func (n *DgDir) ModelCount() int {
	total := len(n.Models)
	for _, child := range n.Children {
		total += child.ModelCount()
	}
	return total
}

func (n *DgDir) prettyPrint(indent string) {
	fmt.Println(indent + n.Name)

	if len(n.Models) > 0 {
		for k := range n.Models {
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
