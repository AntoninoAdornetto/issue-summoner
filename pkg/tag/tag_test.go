package tag_test

import (
	"os"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/stretchr/testify/require"
)

func TestScanForTags_SingleLineCommentEmpty(t *testing.T) {
	file, fileInfo := setup(t,
		`func main(){
		// This is a comment with no annotation.
		// Thus, there should not be a tag returned
		return 0
		}`,
	)

	defer tearDown(file)

	tm := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	ptm := tag.PendedTagManager{TagManager: tm}

	tags, err := ptm.ScanForTags(
		tag.ScanForTagsParams{Path: file.Name(), File: file, FileInfo: fileInfo},
	)

	require.NoError(t, err)
	require.Empty(t, tags)
}

func TestScanForTags_SingleLineCommentOne(t *testing.T) {
	file, fileInfo := setup(t,
		`package main

		import "fmt"

		func main() {
			// @TODO Add Game Loop
			return 0
		}`,
	)

	defer tearDown(file)

	tm := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	ptm := tag.PendedTagManager{TagManager: tm}

	tags, err := ptm.ScanForTags(
		tag.ScanForTagsParams{Path: file.Name(), File: file, FileInfo: fileInfo},
	)

	tag := tags[0]

	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, tag.FileInfo, fileInfo)
	require.Equal(t, tag.Description, "Add Game Loop")
	require.Equal(t, tag.StartLineNumber, uint64(6))
	require.Equal(t, tag.EndLineNumber, uint64(6))
	require.Equal(t, tag.AnnotationLineNum, uint64(6))
}

func TestScanForTags_SingleLineCommentMultiple(t *testing.T) {
	file, fileInfo := setup(t,
		`package main

		import "fmt"

		func main() {
			// @TODO Add feature X
			// Feature X is ...


			// @TODO Add feature Y
			// Feature Y is ...
			return 0
		}`,
	)

	defer tearDown(file)

	tm := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	ptm := tag.PendedTagManager{TagManager: tm}

	tags, err := ptm.ScanForTags(
		tag.ScanForTagsParams{Path: file.Name(), File: file, FileInfo: fileInfo},
	)

	expected := []tag.Tag{
		{
			AnnotationLineNum: uint64(6),
			StartLineNumber:   uint64(6),
			EndLineNumber:     uint64(7),
			Description:       "Add feature X\nFeature X is ...",
			FileInfo:          fileInfo,
		},
		{
			AnnotationLineNum: uint64(10),
			StartLineNumber:   uint64(10),
			EndLineNumber:     uint64(11),
			Description:       "Add feature Y\nFeature Y is ...",
			FileInfo:          fileInfo,
		},
	}

	require.NoError(t, err)
	require.Len(t, tags, 2)
	require.Equal(t, expected, tags)
}

func setup(t *testing.T, fileText string) (*os.File, os.FileInfo) {
	file, err := os.CreateTemp("", "*.go")
	require.NoError(t, err)

	_, err = file.WriteString(fileText)
	require.NoError(t, err)

	err = file.Sync() // Check if this is needed
	require.NoError(t, err)

	_, err = file.Seek(0, 0)
	require.NoError(t, err)

	fileInfo, err := file.Stat()
	require.NoError(t, err)

	return file, fileInfo
}

func tearDown(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}
