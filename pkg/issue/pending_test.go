package issue_test

import (
	"io"
	"os"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

// should handle parsing multiple different types of comments that contain
// annotations in a c file.
func TestScanSingleLineCommentsGo(t *testing.T) {
	srcFile, err := os.Open("../../testdata/test.c")
	require.NoError(t, err)

	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	src, err := io.ReadAll(srcFile)
	require.NoError(t, err)

	err = im.Scan(src, "../../testdata/test.c")
	require.NoError(t, err)

	// see test.c in the testdata directory for how the file looks prior
	// to scanning. The file contains c source code with comments that
	// are parsed. All comments that are annotated with @TEST_TODO should
	// appear in the expected slice.
	expected := []issue.Issue{
		{
			ID:          "test.c-62:95",
			Title:       "inline comment #1",
			Description: "",
			LineNumber:  5,
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			StartIndex:  62,
			EndIndex:    95,
		},
		{
			ID:          "test.c-115:148",
			Title:       "inline comment #2",
			Description: "",
			LineNumber:  6,
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			StartIndex:  115,
			EndIndex:    148,
		},
		{
			ID:          "test.c-192:252",
			Title:       "decode the message and clean up after yourself!",
			Description: "",
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			LineNumber:  10,
			StartIndex:  192,
			EndIndex:    252,
		},
		{
			ID:          "test.c-269:561",
			Title:       "drop a star if you know about this code wars challenge",
			Description: "Digital Cypher assigns to each letter of the alphabet unique number. Instead of letters in encrypted word we write the corresponding number Then we add to each obtained digit consecutive digits from the key",
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			LineNumber:  19,
			StartIndex:  269,
			EndIndex:    561,
		},
	}

	actual := im.GetIssues()
	require.Equal(t, expected, actual)
}

// should walk the testdata directory and return a count that is equal
// to the number of times that the function walk is called.
func TestWalkCountScans(t *testing.T) {
	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	// Walk is called 4 times in total within the testdata directory
	// 1 time for the .gitignore file
	// 1 time for the ignore.sh file
	// 1 time for test.c
	// 1 time for test.log
	// this test does not include gitignore exclude rules. see next test for that.
	expected := 4
	actual, err := im.Walk("../../testdata/")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should walk the testdata directory and return a count that is equal
// to the number of times walk calls the scan method. This time,
// we will add an exclude pattern to the .gitignore file in testdata/.gitignore
func TestWalkCountStepsWithIgnorePatterns(t *testing.T) {
	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	ignoreFile, err := os.OpenFile("../../testdata/.gitignore", os.O_RDWR, 0644)
	require.NoError(t, err)
	defer ignoreFile.Close()

	originalBytes, err := io.ReadAll(ignoreFile)
	require.NoError(t, err)

	ignorePatterns := `
	# Wildcard pattern
	*.log

	# ignore everything within the exclude directory
	exclude/
	`

	_, err = ignoreFile.Seek(0, 0)
	require.NoError(t, err)
	err = ignoreFile.Truncate(0)
	require.NoError(t, err)

	_, err = ignoreFile.WriteString(ignorePatterns)
	require.NoError(t, err)
	err = ignoreFile.Sync()
	require.NoError(t, err)

	defer func() {
		_, err := ignoreFile.Seek(0, 0)
		require.NoError(t, err)
		err = ignoreFile.Truncate(0)
		require.NoError(t, err)
		_, err = ignoreFile.Write(originalBytes)
		if err != nil {
			t.Fatalf("Failed to restore .gitignore contents: %s", err.Error())
		}
		err = ignoreFile.Sync()
		require.NoError(t, err)
	}()

	// walk should call scan one time for test.c and one test for .gitignore
	// might change calling it on .gitignore file but for now it's ok
	expected := 2
	actual, err := im.Walk("../../testdata")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return an error when the root path does not exist
func TestWalkNoneExistentRoot(t *testing.T) {
	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	_, err = im.Walk("unknown-path")
	require.Error(t, err)
}
