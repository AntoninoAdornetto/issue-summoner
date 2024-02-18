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

	tags, err := ptm.ScanForTags(file.Name(), file, fileInfo)

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

	tags, err := ptm.ScanForTags(file.Name(), file, fileInfo)

	tag := tags[0]

	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, fileInfo, tag.FileInfo)
	require.Equal(t, "Add Game Loop", tag.Title)
	require.Equal(t, "", tag.Description)
	require.Equal(t, uint64(6), tag.StartLineNumber)
	require.Equal(t, uint64(6), tag.EndLineNumber)
	require.Equal(t, uint64(6), tag.AnnotationLineNum)
	require.True(t, tag.IsSingleLine)
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

	tags, err := ptm.ScanForTags(file.Name(), file, fileInfo)

	expected := []tag.Tag{
		{
			AnnotationLineNum: uint64(6),
			StartLineNumber:   uint64(6),
			EndLineNumber:     uint64(7),
			Title:             "Add feature X",
			Description:       "Feature X is ...",
			FileInfo:          fileInfo,
			IsSingleLine:      true,
			IsMultiLine:       false,
		},
		{
			AnnotationLineNum: uint64(10),
			StartLineNumber:   uint64(10),
			EndLineNumber:     uint64(11),
			Title:             "Add feature Y",
			Description:       "Feature Y is ...",
			FileInfo:          fileInfo,
			IsSingleLine:      true,
			IsMultiLine:       false,
		},
	}

	require.NoError(t, err)
	require.Len(t, tags, 2)
	require.Equal(t, expected, tags)
}

func TestScanForTags_MultiLineCommentEmpty(t *testing.T) {
	file, fileInfo := setup(t,
		`func main(){
		/*
		This is a multi line comment with no annotation.
		Thus, there should not be a tag struct returned
		*/
		return 0
		}`,
	)

	defer tearDown(file)

	tm := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	ptm := tag.PendedTagManager{TagManager: tm}

	tags, err := ptm.ScanForTags(file.Name(), file, fileInfo)

	require.NoError(t, err)
	require.Empty(t, tags)
}

func TestScanForTags_MultiLineCommentOne(t *testing.T) {
	file, fileInfo := setup(t,
		`func main(){
		/*
		@TODO Add feature X
		This is a multi line comment with a single annotation.
		And some additional information
		*/
		return 0
		}`,
	)

	defer tearDown(file)

	tm := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	ptm := tag.PendedTagManager{TagManager: tm}

	tags, err := ptm.ScanForTags(file.Name(), file, fileInfo)

	tag := tags[0]

	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, fileInfo, tag.FileInfo)
	require.Equal(t, "Add feature X", tag.Title)
	require.Equal(t,
		"This is a multi line comment with a single annotation. And some additional information",
		tag.Description,
	)
	require.Equal(t, uint64(2), tag.StartLineNumber)
	require.Equal(t, uint64(6), tag.EndLineNumber)
	require.Equal(t, uint64(3), tag.AnnotationLineNum)
	require.True(t, tag.IsMultiLine)
	require.False(t, tag.IsSingleLine)
}

func TestScanForTags_MultiLineCommentMany(t *testing.T) {
	file, fileInfo := setup(t,
		`
		/*
		@Author Antonino Adornetto
		@TODO Add feature Y
		Feature Y is ...
		*/

		func main(){
		/*
		@TODO Add feature X
		This is a multi line comment with a single annotation.
		And some additional information
		*/
		return 0
		}`,
	)

	defer tearDown(file)

	tm := tag.TagManager{
		TagName: "@TODO",
		Mode:    "P", // Pending
	}

	ptm := tag.PendedTagManager{TagManager: tm}

	tags, err := ptm.ScanForTags(file.Name(), file, fileInfo)

	expected := []tag.Tag{
		{
			AnnotationLineNum: uint64(4),
			StartLineNumber:   uint64(2),
			EndLineNumber:     uint64(6),
			Title:             "Add feature Y",
			Description:       "@Author Antonino Adornetto Feature Y is ...",
			FileInfo:          fileInfo,
			IsSingleLine:      false,
			IsMultiLine:       true,
		},
		{
			AnnotationLineNum: uint64(10),
			StartLineNumber:   uint64(9),
			EndLineNumber:     uint64(13),
			Title:             "Add feature X",
			Description:       "This is a multi line comment with a single annotation. And some additional information",
			FileInfo:          fileInfo,
			IsSingleLine:      false,
			IsMultiLine:       true,
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
