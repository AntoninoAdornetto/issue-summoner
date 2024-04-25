package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

// should get the correct single/multi line comment notation for a file
// that uses c
func TestNewCommentNotationC(t *testing.T) {
	notation := issue.NewCommentNotation(".c")
	require.Equal(t, notation.SingleLinePrefixRe, issue.CommentNotations[".c"].SingleLinePrefixRe)
	require.Equal(t, notation.MultiLinePrefixRe, issue.CommentNotations[".c"].MultiLinePrefixRe)
	require.Equal(t, notation.MultiLineSuffixRe, issue.CommentNotations[".c"].MultiLineSuffixRe)
	require.Equal(t, notation.NewLinePrefixRe, issue.CommentNotations[".c"].NewLinePrefixRe)
}

// should get the correct single/multi line comment notation for a file
// that uses python
func TestNewCommentNotationPython(t *testing.T) {
	notation := issue.NewCommentNotation(".py")
	require.Equal(t, notation.SingleLinePrefixRe, issue.CommentNotations[".py"].SingleLinePrefixRe)
	require.Equal(t, notation.MultiLinePrefixRe, issue.CommentNotations[".py"].MultiLinePrefixRe)
	require.Equal(t, notation.MultiLineSuffixRe, issue.CommentNotations[".py"].MultiLineSuffixRe)
	require.Equal(t, notation.NewLinePrefixRe, issue.CommentNotations[".py"].NewLinePrefixRe)
}

// should get the correct single/multi line comment notation for a file
// that uses markdown
func TestNewCommentNotationMarkdown(t *testing.T) {
	notation := issue.NewCommentNotation(".md")
	require.Equal(t, notation.SingleLinePrefixRe, issue.CommentNotations[".md"].SingleLinePrefixRe)
	require.Equal(t, notation.MultiLinePrefixRe, issue.CommentNotations[".md"].MultiLinePrefixRe)
	require.Equal(t, notation.MultiLineSuffixRe, issue.CommentNotations[".md"].MultiLineSuffixRe)
	require.Equal(t, notation.NewLinePrefixRe, issue.CommentNotations[".md"].NewLinePrefixRe)
}

// should fall back to a comment notation that uses # to denote comments when the
// comment notations map does not have an key/val store for the file type.
func TestNewCommentNotationDefault(t *testing.T) {
	notation := issue.NewCommentNotation(".sh")
	require.Equal(
		t,
		notation.SingleLinePrefixRe,
		issue.CommentNotations["default"].SingleLinePrefixRe,
	)
	require.Equal(
		t,
		notation.MultiLinePrefixRe,
		issue.CommentNotations["default"].MultiLinePrefixRe,
	)
	require.Equal(
		t,
		notation.MultiLineSuffixRe,
		issue.CommentNotations["default"].MultiLineSuffixRe,
	)
	require.Equal(t, notation.NewLinePrefixRe, issue.CommentNotations["default"].NewLinePrefixRe)
}
