package issue_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

const (
	c_single_line_comment      = "// @TEST_TODO single line comment in c"
	c_single_line_comment_src  = "int main() {return 0;} // @TEST_TODO single line comment w/ src code"
	py_single_line_comment     = "# @TEST_TODO single line comment in python"
	py_single_line_comment_src = "n = 5 # @TEST_TODO single line comment w/ src code"
	md_single_line_comment     = "<!-- @TEST_TODO single line comment in markdown -->"
)

func setupSingleLineComment(ext string, line []byte) issue.SingleLineComment {
	notation := issue.NewCommentNotation(ext)
	prefixLocations, _ := notation.FindPrefixLocations(line)
	buf := bytes.NewBuffer(line)
	scanner := bufio.NewScanner(buf)
	scanner.Scan()

	return issue.SingleLineComment{
		Annotation:               annotation,
		AnnotationIndicator:      false,
		FileName:                 fmt.Sprintf("test%s", ext),
		FilePath:                 fmt.Sprintf("/home/user/app/test%s", ext),
		PrefixRe:                 issue.CommentNotations[ext].SingleLinePrefixRe,
		SuffixRe:                 issue.CommentNotations[ext].SingleLineSuffixRe,
		CommentNotationLocations: prefixLocations,
		Scanner:                  scanner,
	}
}

// should parse a single line comment and remove the comment symbols and the annotation
// and return a Comment object. Note: single line comment will not have a description.
func TestParseCommentC(t *testing.T) {
	slc := setupSingleLineComment(".c", []byte(c_single_line_comment))

	expected := []issue.Comment{
		{
			ID:                   "test.c-1",
			Title:                "single line comment in c",
			Description:          "",
			FileName:             "test.c",
			FilePath:             "/home/user/app/test.c",
			StartLineNumber:      1,
			EndLineNumber:        1,
			AnnotationLineNumber: 1,
			ColumnLocations:      [][]int{{2}},
		},
	}

	actual, err := slc.ParseComment(1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return an int slice with 1 item in it that signifies where we should begin
// to slice. Slicing allows us to strip the comment off the resulting string so that it
// is not part of the comment title. In C, there is no suffix to conclude a single line
// comment. The end of the line would deem it as complete. Thus, the resulting int slice
// should contain only 1 integer.

// should return an int slice containing the index of where we should begin to slice
// the byte array so that the comment symbols can be removed and used to build an issue/comment.
// In C, there is no suffix used to conclude a single line comment. The end of the line will deem
// it as complete.
func TestFindCutLocationsC(t *testing.T) {
	line := []byte(c_single_line_comment)
	slc := setupSingleLineComment(".c", line)
	expected := []int{2} // last comment symbol is located at index 1
	actual, err := slc.FindCutLocations(line)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should also work when the comment symbols are not at the begining of the line
// but further into the slice (line). We may see this in scenarios where there
// is some amount of source code, and after the source code tokens, a comment
// is inserted. This test should assert that we still receive a proper index
// of where we can cut the comment out of the byte slice (line).
func TestFindCutLocationsCWithSourceCode(t *testing.T) {
	line := []byte(c_single_line_comment_src)
	slc := setupSingleLineComment(".c", line)
	expected := []int{25} // last byte of the comment is at index 24
	actual, err := slc.FindCutLocations(line)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return an int slice containing the index of where we should begin to slice
// the byte array for python.
func TestFindCutLocationsPython(t *testing.T) {
	line := []byte(py_single_line_comment)
	slc := setupSingleLineComment(".py", line)
	expected := []int{1} // the comment symbol is located at index 0
	actual, err := slc.FindCutLocations(line)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return an int slice containing the index of where we should begin to slice
// the byte array for python when the line starts with source code.
func TestFindCutLocationsPythonWithSourceCode(t *testing.T) {
	line := []byte(py_single_line_comment_src)
	slc := setupSingleLineComment(".py", line)
	expected := []int{7} // the comment symbol is located at index 6
	actual, err := slc.FindCutLocations(line)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return a nil int slice when working with c because
// there is no closing comment notation (suffix)
func TestFindSuffixLocationsC(t *testing.T) {
	line := []byte(c_single_line_comment)
	slc := setupSingleLineComment(".c", line)
	expected := []int(nil)
	actual := slc.FindSuffixLocations(line)
	require.Equal(t, expected, actual)
}

// should return a nil int slice when working with python because
// there is no closing comment notation (suffix)
func TestFindSuffixLocationsPython(t *testing.T) {
	line := []byte(py_single_line_comment)
	slc := setupSingleLineComment(".py", line)
	expected := []int(nil)
	actual := slc.FindSuffixLocations(line)
	require.Equal(t, expected, actual)
}

// should return an int slice with a len of 2 for markdown since it
// does make use of a suffix to conclude a single line comment.
// the int slice will contain the index of where the suffix starts and where it ends
func TestFindSuffixLocationsMarkdown(t *testing.T) {
	line := []byte(md_single_line_comment)
	slc := setupSingleLineComment(".md", line)
	expected := []int{48, 51}
	actual := slc.FindSuffixLocations(line)
	require.Equal(t, expected, actual)
}

// should take the indices returned from the FindCutLocations func
// and use them to cut the comment symbols out of the byte slice (line).
// the result is that we get a string that does not denote any single line
// comments
func TestSliceC(t *testing.T) {
	line := []byte(c_single_line_comment)
	slc := setupSingleLineComment(".c", line)
	locations, err := slc.FindCutLocations(line)
	require.NoError(t, err)
	expected := []byte("@TEST_TODO single line comment in c") // comment symbols are removed
	actual, err := slc.Slice(line, locations)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should take the indices returned from the FindCutLocations func
// and use them to cut the comment symbols out of the byte slice (line).
// this time we are adding c source code before the comment is denoted.
func TestSliceCWithSourceCode(t *testing.T) {
	line := []byte(c_single_line_comment_src)
	slc := setupSingleLineComment(".c", line)
	locations, err := slc.FindCutLocations(line)
	require.NoError(t, err)
	expected := []byte("@TEST_TODO single line comment w/ src code") // comment symbols are removed
	actual, err := slc.Slice(line, locations)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

}
