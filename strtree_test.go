package dsa_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipsusila/dsa"
	"github.com/stretchr/testify/assert"
)

func loadTree(t *testing.T, filename string, arg dsa.StringTreeArg) *dsa.Tree[string] {
	fd, err := os.Open(filename)
	if !assert.NoError(t, err, "Open file shall not error") {
		assert.Fail(t, "Failed to open file")
	}
	defer fd.Close()

	outDir := filepath.Dir(filename)
	extName := filepath.Ext(filename)
	baseName := strings.ReplaceAll(filepath.Base(filename), extName, "")
	outName := filepath.Join(outDir, baseName+"_out"+extName)
	fo, err := os.Create(outName)
	if !assert.NoError(t, err, "Open file shall not error") {
		assert.Fail(t, "Failed to create file")
	}
	defer fo.Close()
	tree := dsa.NewStringTree(fd, arg)
	tree.PrintTo(fo)

	return tree
}

func TestTree(t *testing.T) {
	var files = []string{
		filepath.Join("_testdata", "regency.txt"),
		filepath.Join("_testdata", "city.txt"),
	}
	for idx, filename := range files {
		arg := dsa.StringTreeArg{
			Name: filename,
		}
		tree := loadTree(t, filename, arg)

		var values []string
		var nd *dsa.Node[string]
		var ok bool

		if idx == 0 {
			// Entries:
			// 1. KAB. OGAN KOMERING ULU
			// 2. KAB. OGAN KOMERING ULU TIMUR
			// 3. KAB. OGAN KOMERING ULU SELATAN
			values = []string{"KAB.", "OGAN", "KOMERING", "ULU"}
			nd, ok = tree.Match(values)
			assert.True(t, ok, "Should match partially")
			assert.NotEmpty(t, nd.Nodes, "Node shall has children")

			nd, ok = tree.ExactMatch(values)
			assert.True(t, ok, "Should match exactly (with children)")
			assert.NotEmpty(t, nd.Nodes, "Leaf shall has children")
		} else if idx == 1 {
			values = []string{"KOTA", "ADM.", "JAKARTA", "SELATAN"}
			nd, ok = tree.ExactMatch(values)
			assert.True(t, ok, "Should match exactly")
			assert.Empty(t, nd.Nodes, "Node shall not has children")
		}
	}
}

func TestTreePrefix(t *testing.T) {
	var files = []string{
		filepath.Join("_testdata", "regency.txt"),
		filepath.Join("_testdata", "city.txt"),
	}
	exPrefixs := []string{"ADM.", "KAB", "KAB.", "KOTA"}
	arg := dsa.StringTreeArg{
		Name: "RegencyFilter",
		FnFilter: func(v string) bool {
			for _, prefix := range exPrefixs {
				if prefix == v {
					return true
				}
			}
			return false
		},
		FnTransform: func(v string) string {
			switch v {
			case "KEP.":
				return "KEPULAUAN"
			}
			return v
		},
	}
	for idx, filename := range files {
		tree := loadTree(t, filename, arg)

		var values []string
		var nd *dsa.Node[string]
		var ok bool

		if idx == 0 {
			// Entries:
			// 1. KAB. OGAN KOMERING ULU
			// 2. KAB. OGAN KOMERING ULU TIMUR
			// 3. KAB. OGAN KOMERING ULU SELATAN
			values = []string{"OGAN", "KOMERING", "ULU"}
			nd, ok = tree.Match(values)
			assert.True(t, ok, "Should match partially")
			assert.NotEmpty(t, nd.Nodes, "Node shall has children")

			nd, ok = tree.ExactMatch(values)
			assert.True(t, ok, "Should match exactly (with children)")
			assert.NotEmpty(t, nd.Nodes, "Leaf shall has children")
		} else if idx == 1 {
			values = []string{"JAKARTA", "SELATAN"}
			nd, ok = tree.ExactMatch(values)
			assert.True(t, ok, "Should match exactly")
			assert.Empty(t, nd.Nodes, "Node shall not has children")
		}
	}
}
