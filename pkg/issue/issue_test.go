package issue

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testAnnotation = []byte("@TEST_ANNOTATION")
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
