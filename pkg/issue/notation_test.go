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

// should find the start and end indices of the single line prefix bytes and return the comment
// type as a single line comment when provided with a byte slice that contains a sinle line comment in c
func TestFindPrefixLocationSingleLineCommentC(t *testing.T) {
	notation := issue.NewCommentNotation(".c")
	line := []byte("int n = 5; // single line comment using the c programming language")
	expectedLocations, expectedCommentType := []int{11, 13}, issue.SINGLE_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should find the start and end indices of the multi line prefix bytes and return the comment
// type as a multi line comment when provided with a byte slice that contains a sinle line comment in c
func TestFindPrefixLocationMultiLineCommentC(t *testing.T) {
	notation := issue.NewCommentNotation(".c")
	line := []byte("/* multi line comment using the c programming language")
	expectedLocations, expectedCommentType := []int{0, 2}, issue.MULTI_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should return an empty integer slice and an empty string as the comment type when a byte slice
// is provided as input that does not contain a single or multi line comment
func TestFindPrefixLocationNoCommentC(t *testing.T) {
	notation := issue.NewCommentNotation(".c")
	line := []byte("int main(argc int argv **char) {return 0; }")
	expectedLocations, expectedCommentType := []int(nil), ""
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should find the start and end indices of the single line prefix bytes and return the comment
// type as a single line comment when provided with a byte slice that contains a sinle line comment in go
func TestFindPrefixLocationSingleLineCommentGo(t *testing.T) {
	notation := issue.NewCommentNotation(".go")
	line := []byte("// single line comment using the go programming language")
	expectedLocations, expectedCommentType := []int{0, 2}, issue.SINGLE_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should find the start and end indices of the multi line prefix bytes and return the comment
// type as a multi line comment when provided with a byte slice that contains a sinle line comment in go
func TestFindPrefixLocationMultiLineCommentGo(t *testing.T) {
	notation := issue.NewCommentNotation(".go")
	line := []byte("/* multi line comment using the go programming language")
	expectedLocations, expectedCommentType := []int{0, 2}, issue.MULTI_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should return an empty integer slice and an empty string as the comment type when a byte slice
// is provided as input that does not contain a single or multi line comment
func TestFindPrefixLocationNoCommentGo(t *testing.T) {
	notation := issue.NewCommentNotation(".go")
	line := []byte("func main() { fmt.Println('Hello World') }")
	expectedLocations, expectedCommentType := []int(nil), ""
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should find the start and end indices of the single line prefix bytes and return the comment
// type as a single line comment when provided with a byte slice that contains a sinle line comment in python
func TestFindPrefixLocationSingleLineCommentPython(t *testing.T) {
	notation := issue.NewCommentNotation(".py")
	line := []byte("# single line comment using the python programming language")
	expectedLocations, expectedCommentType := []int{0, 1}, issue.SINGLE_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should find the start and end indices of the multi line prefix bytes and return the comment
// type as a multi line comment when provided with a byte slice that contains a sinle line comment in c
func TestFindPrefixLocationMultiLineCommentPython(t *testing.T) {
	notation := issue.NewCommentNotation(".py")
	line := []byte("\"\"\" multi line comment using the py programming language \"\"\"")
	expectedLocations, expectedCommentType := []int{0, 3}, issue.MULTI_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should return an empty integer slice and an empty string as the comment type when a byte slice
// is provided as input that does not contain a single or multi line comment
func TestFindPrefixLocationNoCommentPython(t *testing.T) {
	notation := issue.NewCommentNotation(".py")
	line := []byte("n = 5")
	expectedLocations, expectedCommentType := []int(nil), ""
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}

// should find the start and end indices of the single line prefix bytes and return the comment
// type as a single line comment when provided with a byte slice that contains a sinle line comment in markdown
func TestFindPrefixLocationSingleLineCommentMarkdown(t *testing.T) {
	notation := issue.NewCommentNotation(".md")
	line := []byte("<!-- single line comment using markdown -->")
	expectedLocations, expectedCommentType := []int{0, 4}, issue.SINGLE_LINE_COMMENT
	actualLocations, actualCommentType := notation.FindPrefixLocations(line)
	require.Equal(t, expectedLocations, actualLocations)
	require.Equal(t, expectedCommentType, actualCommentType)
}
