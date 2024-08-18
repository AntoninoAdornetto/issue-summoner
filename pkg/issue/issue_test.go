package issue

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	// currentIssueCount is the actual number of annotations that this project currently has
	currentIssueCount = 4
	testAnnotation    = []byte("@TEST_ANNOTATION")
)

func TestNewIssueManager(t *testing.T) {
	testCases := []struct {
		name     string
		expected *IssueManager
		mode     IssueMode
		err      bool
	}{
		{
			name: "Create a new issue manager with the correct mode and annotation fields when invoked with scan mode",
			expected: &IssueManager{
				annotation: testAnnotation,
				Issues:     []Issue{},
				mode:       IssueModeScan,
				os:         runtime.GOOS,
			},
			mode: IssueModeScan,
			err:  false,
		},
		{
			name: "Create a new issue manager with the correct mode and annotation when invoked with purge mode",
			expected: &IssueManager{
				// when purging comments, the annotation is constructed in a way that will allow the lexer package
				// to discover annotations that have an issue id, enclosed within parans, appended to the annotation.
				annotation: []byte("@TEST_ANNOTATION\\(\\d+\\)"),
				Issues:     []Issue{},
				mode:       IssueModePurge,
				os:         runtime.GOOS,
			},
			mode: IssueModePurge,
			err:  false,
		},
		{
			name: "Create a new issue manager with the correct mode and annotation fields when invoked with report mode",
			expected: &IssueManager{
				annotation: testAnnotation,
				Issues:     []Issue{},
				mode:       IssueModeReport,
				os:         runtime.GOOS,
			},
			mode: IssueModeReport,
			err:  false,
		},
		{
			name:     "Returns an error when provided a mode that is not supported",
			expected: nil,
			mode:     "unsupported-mode",
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := NewIssueManager(testAnnotation, tc.mode)

			if tc.err {
				require.Error(t, err)
				require.Nil(t, manager)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, manager)

				if tc.mode == IssueModePurge {
					require.Equal(t, `@TEST_ANNOTATION\(\d+\)`, string(tc.expected.annotation))
				}
			}
		})
	}
}

func TestWalk(t *testing.T) {
	manager, err := NewIssueManager([]byte("@TODO"), IssueModeScan)
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	err = manager.Walk(filepath.Join(wd, "../../"))
	require.NoError(t, err)

	require.Len(t, manager.Issues, currentIssueCount)
}

func BenchmarkWalk(b *testing.B) {
	manager, err := NewIssueManager([]byte("@TODO"), IssueModeScan)
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
