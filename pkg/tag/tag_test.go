package tag_test

import (
	"os"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/stretchr/testify/require"
)

func TestFindPendingTags(t *testing.T) {
	file, err := os.CreateTemp("", "*.gitignore")
	require.NoError(t, err)

	defer os.Remove(file.Name())
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

	tagManager := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	mockIgnoreFileOpener := MockIgnoreFileOpener{File: file}
	pendedTagManager := tag.PendedTagManager{TagManager: tagManager}

	pendedTags, err := pendedTagManager.FindTags(file.Name(), mockIgnoreFileOpener)
	require.NoError(t, err)
	require.Len(t, pendedTags, 1)
	require.Equal(t, uint64(9), pendedTags[0].LineNum)
}
