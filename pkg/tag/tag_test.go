package tag_test

import (
	"os"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/stretchr/testify/require"
)

func TestFindPendingTags_SingleLineComment(t *testing.T) {
	file, err := os.CreateTemp("", "*.go")
	require.NoError(t, err)

	defer os.Remove(file.Name())
	defer file.Close()

	// Will test more languages, starting with `Go` for now
	_, err = file.WriteString(`
		package main

		import "fmt"

		func main() {
			fmt.Printf("Hello World\n")

			// @TODO - Add Game Loop
		}
	`)
	require.NoError(t, err)

	err = file.Sync()
	require.NoError(t, err)

	_, err = file.Seek(0, 0)
	require.NoError(t, err)

	fileInfo, err := file.Stat()
	require.NoError(t, err)

	tagManager := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}
	pendedTagManager := tag.PendedTagManager{TagManager: tagManager}

	pendedTags, err := pendedTagManager.ScanForTags(tag.ScanForTagsParams{
		Path:     file.Name(),
		File:     file,
		FileInfo: fileInfo,
	})
	require.NoError(t, err)
	require.Len(t, pendedTags, 1)
	require.Equal(t, uint64(9), pendedTags[0].LineNumber)
}
