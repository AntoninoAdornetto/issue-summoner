package issue_test

import (
	"os"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestScan_SingleLine1Item(t *testing.T) {
	file, _ := setup(
		t,
		"func main() {\n\t// @TEST_TODO Test Me\n\t// Write test cases that bring up the code coverage\n}\n",
	)

	defer tearDown(file)

	pi, err := issue.GetIssueManager(issue.PENDING_ISSUE, "@TEST_TODO")
	require.NoError(t, err)

	actual, err := pi.Scan(file)
	require.NoError(t, err)

	expected := []issue.Issue{
		{
			AnnotationLineNumber: 2,
			StartLineNumber:      2,
			EndLineNumber:        3,
			Title:                "Test Me",
			Description:          "Write test cases that bring up the code coverage",
			IsSingleLine:         true,
			IsMultiLine:          false,
		},
	}

	require.Equal(t, expected, actual)
}

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
