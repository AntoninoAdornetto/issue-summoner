package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedCommentManager(t *testing.T) {
	arg := issue.NewCommentManagerParams{
		Annotation:      annotation,
		CommentType:     "",
		FilePath:        "",
		FileName:        "",
		StartLineNumber: 1,
		Locations:       []int{},
		PrefixRe:        nil,
		SuffixRe:        nil,
		Scanner:         nil,
	}

	manager, err := issue.NewCommentManager(arg)
	require.Error(t, err)
	require.Nil(t, manager)
}

func TestSingleLineCommentManager(t *testing.T) {
	arg := issue.NewCommentManagerParams{
		Annotation:      annotation,
		CommentType:     issue.SINGLE_LINE_COMMENT,
		FilePath:        "",
		FileName:        "",
		StartLineNumber: 1,
		Locations:       []int{0, 3},
		PrefixRe:        issue.CommentNotations[".c"].SingleLinePrefixRe,
		SuffixRe:        nil,
		Scanner:         nil,
	}

	manager, err := issue.NewCommentManager(arg)
	require.NoError(t, err)
	require.IsType(t, &issue.SingleLineComment{}, manager)
}
