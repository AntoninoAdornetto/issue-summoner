package issue_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestScan_SingleLine1Item(t *testing.T) {
	file, fileInfo := setup(
		t,
		`func main(){
			// @TEST_TODO Test Me
			// Write test cases that bring up the code coverage
		}`,
	)

	defer tearDown(file)

	pi, err := issue.GetIssueManager(issue.PENDING_ISSUE, "@TEST_TODO")
	require.NoError(t, err)

	err = pi.Scan(file)
	require.NoError(t, err)
	actual := pi.GetIssues()

	expected := []issue.Issue{
		{
			AnnotationLineNumber: 2,
			StartLineNumber:      2,
			EndLineNumber:        3,
			Title:                "Test Me",
			Description:          "Write test cases that bring up the code coverage",
			IsSingleLine:         true,
			IsMultiLine:          false,
			FileInfo:             fileInfo,
			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 2),
		},
	}

	require.Equal(t, expected, actual)
}

func TestScan_SingleLineMultipleItems(t *testing.T) {
	file, fileInfo := setup(
		t,
		`func main() {
		// @TEST_TODO Test Me

		fmt.Printf("hello world\n")
		// @TEST_TODO Test Me as well
		}`,
	)

	defer tearDown(file)

	pi, err := issue.GetIssueManager(issue.PENDING_ISSUE, "@TEST_TODO")
	require.NoError(t, err)

	err = pi.Scan(file)
	require.NoError(t, err)
	actual := pi.GetIssues()

	expected := []issue.Issue{
		{
			AnnotationLineNumber: 2,
			StartLineNumber:      2,
			EndLineNumber:        2,
			Title:                "Test Me",
			Description:          "",
			IsSingleLine:         true,
			IsMultiLine:          false,
			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 2),
			FileInfo:             fileInfo,
		},
		{
			AnnotationLineNumber: 5,
			StartLineNumber:      5,
			EndLineNumber:        5,
			Title:                "Test Me as well",
			Description:          "",
			IsSingleLine:         true,
			IsMultiLine:          false,
			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 5),
			FileInfo:             fileInfo,
		},
	}

	require.Equal(t, expected, actual)
}

func TestScanMultiCommentLine1Item(t *testing.T) {
	file, fileInfo := setup(
		t,
		`func main(){
			/* @TEST_TODO Test Me
				Write test cases that bring up the code coverage
	*/
		}`,
	)

	defer tearDown(file)

	pi, err := issue.GetIssueManager(issue.PENDING_ISSUE, "@TEST_TODO")
	require.NoError(t, err)

	err = pi.Scan(file)
	require.NoError(t, err)
	actual := pi.GetIssues()

	expected := []issue.Issue{
		{
			AnnotationLineNumber: 2,
			StartLineNumber:      2,
			EndLineNumber:        3,
			Title:                "Test Me",
			Description:          "Write test cases that bring up the code coverage",
			IsSingleLine:         false,
			IsMultiLine:          true,
			FileInfo:             fileInfo,
			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 2),
		},
	}

	require.Equal(t, expected, actual)
}

// func TestScanMultiCommentLineManyItems(t *testing.T) {
// 	file, fileInfo := setup(
// 		t,
// 		`func main(){
// 			/* @TEST_TODO First Multi Comment
// 			First multi comment details
// 			*/
//
//
// 			/* @TEST_TODO Second Multi Comment
// 			Second multi comment details
// 			*/
//
// 			/*
// 			@TEST_TODO Third Multi Comment
// 			Third multi comment details
// 			Span for multiple lines...
// 			*/
// 		}`,
// 	)
//
// 	defer tearDown(file)
//
// 	pi, err := issue.GetIssueManager(issue.PENDING_ISSUE, "@TEST_TODO")
// 	require.NoError(t, err)
//
// 	err = pi.Scan(file)
// 	require.NoError(t, err)
// 	actual := pi.GetIssues()
//
// 	expected := []issue.Issue{
// 		{
// 			AnnotationLineNumber: 2,
// 			StartLineNumber:      2,
// 			EndLineNumber:        4,
// 			Title:                "First Multi Comment",
// 			Description:          "First multi comment details",
// 			IsSingleLine:         false,
// 			IsMultiLine:          true,
// 			FileInfo:             fileInfo,
// 			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 2),
// 		},
// 		{
// 			AnnotationLineNumber: 7,
// 			StartLineNumber:      7,
// 			EndLineNumber:        9,
// 			Title:                "Second Multi Comment",
// 			Description:          "Second multi comment details",
// 			IsSingleLine:         false,
// 			IsMultiLine:          true,
// 			FileInfo:             fileInfo,
// 			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 7),
// 		},
// 		{
// 			AnnotationLineNumber: 12,
// 			StartLineNumber:      11,
// 			EndLineNumber:        15,
// 			Title:                "Third Multi Comment",
// 			Description:          "Third multi comment details Span for multiple lines...",
// 			IsSingleLine:         false,
// 			IsMultiLine:          true,
// 			FileInfo:             fileInfo,
// 			ID:                   fmt.Sprintf("%s-%d", fileInfo.Name(), 11),
// 		},
// 	}
//
// 	require.Equal(t, expected, actual)
// }

func setup(t *testing.T, text string) (*os.File, os.FileInfo) {
	file, err := os.CreateTemp("", "*.go")
	require.NoError(t, err)

	_, err = file.WriteString(text)
	require.NoError(t, err)

	err = file.Sync()
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
