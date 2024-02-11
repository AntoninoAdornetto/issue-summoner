package tag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/stretchr/testify/require"
)

func TestFindTags(t *testing.T) {
	file, err := os.CreateTemp("", "impl.go")
	require.NoError(t, err)

	defer file.Close()

	// Line Number 9 is where the Tag is located
	file.WriteString(`
		package main

		import "fmt"

		func main(){
			fmt.Printf("Hello World\n")

			// @TODO - Add Game Loop
		}
	`)

	_, err = file.Seek(0, 0)
	require.NoError(t, err)

	mockIgnoreFileOpener := MockIgnoreFileOpener{File: file}

	tags, err := tag.FindTags(
		file.Name(),
		tag.CompileTagRegexp("@TODO"),
		mockIgnoreFileOpener,
	)

	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, uint64(9), tags[0].LineNum)
	require.Equal(t, filepath.Ext(file.Name()), tags[0].FileExt)
}
