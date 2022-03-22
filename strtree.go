package dsa

import (
	"bufio"
	"io"
	"strings"
)

// StringTreeArg for generating tree
type StringTreeArg struct {
	Sep         string
	Name        string
	FnFilter    FilterFn[string]
	FnTransform TransformFn[string]
}

// NewStringTree with given name
func NewStringTree(r io.Reader, arg StringTreeArg) *Tree[string] {
	// default parameter
	sep := Default(arg.Sep, " ")
	name := Default(arg.Name, "")
	filter := arg.FnFilter
	transform := arg.FnTransform
	if transform == nil {
		transform = func(v string) string { return v }
	}

	// construct tree
	tree := NewTree[string](name)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), sep)
		n := len(values)
		if filter != nil {
			n = 0
			for i := 0; i < len(values); i++ {
				if filter(values[i]) {
					continue
				}
				values[n] = transform(values[i])
				n++
			}
		}
		tree.Insert(values[:n])
	}
	tree.Sort(Ascending)

	return tree
}
