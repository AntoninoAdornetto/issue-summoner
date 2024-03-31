package issue_test

import (
	"strings"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

// TestParseCommentContents_SingleLine is an assertion that we may
// commonly run into. We have a source code file with a single-line
// comment that does not expand for many lines. By expanding I mean
// Example: // @TEST_TODO write unit tests...
// It's a one line comment with 1 annotation and the comment doesn't
// go on for many lines.
func TestParseCommentContents_SingleLine(t *testing.T) {
	line := "// @TEST_TODO increase test cases for ...."
	expected := "@TEST_TODO increase test cases for ...."

	stack := issue.CommentStack{Items: make([]string, 0)}
	builder := strings.Builder{}
	comment := issue.GetCommentSymbols(".c")

	result, err := comment.ParseCommentContents(line, &builder, stack)
	require.NoError(t, err)
	require.Equal(t, expected, result.String())
}

// TestParseCommentContents_SingleLineExpand asserts that we can
// build ontop of single-line comments. In other words, we may have
// single-line coments that expand for many lines.
// func TestParseCommentContents_SingleLineExpand(t *testing.T) {
// 	line := "// @TEST_TODO increase code coverage\n\t"
// 	line2 := "// it's important that we increase code coverage.\n\t"
// 	line3 := "// it will make our boss happy :)\n\t"
// 	lines := []string{line, line2, line3}
//
// 	expected := `@TEST_TODO increase code coverage
// 	it's important that we increase code coverage.
// 	it will make our boss happy :0
// 	`
//
// 	stack := issue.CommentStack{Items: make([]string, 0)}
// 	builder := strings.Builder{}
// 	comment := issue.GetCommentSymbols(".c")
//
// 	for _, l := range lines {
// 		_, err := comment.ParseCommentContents(l, &builder, stack)
// 		require.NoError(t, err)
// 	}
//
// 	require.Equal(t, expected, builder.String())
// }

func TestCommentSymbols_C(t *testing.T) {
	m := issue.GetCommentSymbols(".c")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_GO(t *testing.T) {
	m := issue.GetCommentSymbols(".go")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_JS(t *testing.T) {
	m := issue.GetCommentSymbols(".js")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_TS(t *testing.T) {
	m := issue.GetCommentSymbols(".ts")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_TSX(t *testing.T) {
	m := issue.GetCommentSymbols(".tsx")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_JSX(t *testing.T) {
	m := issue.GetCommentSymbols(".jsx")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_PYTHON(t *testing.T) {
	m := issue.GetCommentSymbols(".py")
	expectedSingleLineSymbols := []string{"#"}
	expectedMultiLineStartSymbols := []string{"\"\"\"", "'''"}
	expectedMultiLineEndSymbols := []string{"\"\"\"", "'''"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_MARKDOWN(t *testing.T) {
	m := issue.GetCommentSymbols(".md")
	expectedMultiLineStartSymbols := []string{"<!--"}
	expectedMultiLineEndSymbols := []string{"-->"}
	require.Empty(t, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}
