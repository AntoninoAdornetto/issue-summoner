package issue_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

var (
	// currentIssueCount is the current number of issues contained in the entire issue-summoner project.
	// The value will change as issues, contained in this project, are resolved and added
	currentIssueCount = 3
	testAnnotation    = []byte("@TEST_ANNOTATION")
)

func TestNewIssueManager(t *testing.T) {
	testCases := []struct {
		name     string
		mode     issue.IssueMode
		expected *issue.IssueManager
		err      bool
	}{
		{
			name: "Should create a new issue manager when invoked with scan mode",
			mode: issue.IssueModeScan,
			expected: &issue.IssueManager{
				Annotation: testAnnotation,
				Issues:     []issue.Issue{},
				IssueMap:   nil,
			},
			err: false,
		},
		{
			name: "Should create a new issue manager when invoked with purge mode",
			mode: issue.IssueModePurge,
			expected: &issue.IssueManager{
				// when purging comments, the annotation is constructed in a way that will allow the lexer package
				// to discover annotations that have an issue id, enclosed within parans, appended to the annotation.
				Annotation: []byte("@TEST_ANNOTATION\\(#\\d+\\)"),
				Issues:     []issue.Issue{},
				IssueMap:   nil,
			},
			err: false,
		},
		{
			name: "Should create a new issue manager when invoked with report mode",
			mode: issue.IssueModeReport,
			expected: &issue.IssueManager{
				Annotation: testAnnotation,
				Issues:     []issue.Issue{},
				IssueMap:   make(map[string][]issue.IssueMapEntry),
			},
			err: false,
		},
		{
			name:     "Should return an error when invoked with a mode that isn't supported",
			mode:     "unsupported-mode",
			expected: nil,
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := issue.NewIssueManager(testAnnotation, tc.mode)
			if tc.err {
				require.Error(t, err)
				require.Nil(t, manager)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected.Annotation, manager.Annotation)
				require.Equal(t, tc.expected.Issues, manager.Issues)
				require.Equal(t, tc.expected.IssueMap, manager.IssueMap)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	testCases := []struct {
		name       string
		mode       issue.IssueMode
		annotation []byte
		expected   int // expected number of issues after walking the working tree
	}{
		{
			name:       "Should return the correct amount issues for the issue-summoner project when in scan mode",
			mode:       issue.IssueModeScan,
			annotation: []byte("@TODO"),
			expected:   currentIssueCount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := issue.NewIssueManager(tc.annotation, tc.mode)
			require.NoError(t, err)

			wd, err := os.Getwd()
			require.NoError(t, err)

			err = manager.Walk(filepath.Join(wd, "../../"))
			require.NoError(t, err)
			require.Len(t, manager.Issues, tc.expected)
		})
	}
}

func BenchmarkWalk(b *testing.B) {
	manager, err := issue.NewIssueManager([]byte("@TODO"), issue.IssueModeScan)
	if err != nil {
		b.Fatalf("Failed to create IssueManager: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get working directory: %v", err)
	}

	path := filepath.Join(wd, "../../")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		manager.Issues = nil

		err := manager.Walk(path)
		if err != nil {
			b.Fatalf("Walk func failed: %v", err)
		}

		if len(manager.Issues) != currentIssueCount {
			b.Fatalf("Expected %d issues, but got %d", currentIssueCount, len(manager.Issues))
		}
	}
}

func TestWriteIssues(t *testing.T) {
	testCases := []struct {
		name           string
		srcPath        string // src file path
		dstPath        string // file we will write to
		validationPath string // validates if the issues written to outPath match what we expect for the test
		issueCount     int
	}{
		{
			name:           "Should write all issue ids to the annotation locations specified in the inPath src code file",
			srcPath:        "./testdata/go/write_id.go",
			dstPath:        "./testdata/go/write_id_out.txt",
			validationPath: "./testdata/go/write_id_expected.txt",
			issueCount:     3, // we should have 3 reported issue ids writen to each annotation in (write_id.go)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := issue.NewIssueManager(testAnnotation, issue.IssueModeReport)
			require.NoError(t, err)

			paths := []string{tc.srcPath, tc.dstPath, tc.validationPath}
			for i, p := range paths {
				abs, err := filepath.Abs(p)
				require.NoError(t, err)
				paths[i] = abs
			}

			err = manager.Scan(tc.srcPath)
			require.NoError(t, err)
			require.Len(t, manager.Issues, tc.issueCount)

			for index, currentIssue := range manager.Issues {
				manager.IssueMap[tc.dstPath] = append(
					manager.IssueMap[tc.dstPath],
					issue.IssueMapEntry{
						Index:      currentIssue.Index,
						ReportedID: index + 1,
					},
				)
			}

			// assert all 3 issues made it into the IssueMap
			require.Len(t, manager.IssueMap[tc.dstPath], tc.issueCount)

			// Helps create a deterministic test by resetting the dst file bytes
			// to that of the srcFile on each test run
			prepareDstFile(t, tc.srcPath, tc.dstPath)
			err = manager.WriteIssues(tc.dstPath)
			require.NoError(t, err)
			match(t, tc.dstPath, tc.validationPath)
		})
	}
}

func prepareDstFile(t *testing.T, src, dst string) {
	srcFile, err := os.Open(src)
	require.NoError(t, err)
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_RDWR, 0666)
	require.NoError(t, err)
	defer dstFile.Close()

	srcCode, err := io.ReadAll(srcFile)
	require.NoError(t, err)

	err = dstFile.Truncate(0)
	require.NoError(t, err)

	_, err = dstFile.Seek(0, io.SeekStart)
	require.NoError(t, err)

	_, err = dstFile.Write(srcCode)
	require.NoError(t, err)

	err = dstFile.Sync()
	require.NoError(t, err)
}

func match(t *testing.T, dstPath, validationPath string) {
	outFile, err := os.Open(dstPath)
	require.NoError(t, err)
	defer outFile.Close()

	expectedFile, err := os.Open(validationPath)
	require.NoError(t, err)
	defer expectedFile.Close()

	a, err := io.ReadAll(outFile)
	require.NoError(t, err)

	b, err := io.ReadAll(expectedFile)
	require.NoError(t, err)

	matched := bytes.Equal(a, b)
	require.True(t, matched)
}
