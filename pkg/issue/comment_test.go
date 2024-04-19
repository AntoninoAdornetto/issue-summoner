package issue_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

const (
	single_line_prefix_c      = "//"
	single_line_prefix_go     = "//"
	single_line_prefix_python = "#"
	multi_line_prefix_c       = "/*"
	multi_line_prefix_go      = "/*"
	multi_line_prefix_python  = "'''"
	file_ext_c                = ".c"
	file_ext_py               = ".py"
	file_ext_go               = ".go"
)

// should return the prefix index as 0 for a c file
func TestFindPrefixIndexSingleLineInC(t *testing.T) {
	line := fmt.Sprintf("%s single line comment in c", single_line_prefix_c)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return -1 when no single line comment prefix exists
func TestFindPrefixIndexSingleLineNotFoundInC(t *testing.T) {
	line := "int main() {return 0}" // no prefix added
	fields := strings.Fields(line)
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := -1
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// since we didn't find the prefix, the stack should be empty
	require.Len(t, notation.Stack.Items, 0)
}

// should return the prefix index as 4 for a c file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexSingleLineAfterSourceCodeInC(t *testing.T) {
	line := fmt.Sprintf("int main() {return 0}; %s single line comment in c", single_line_prefix_c)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := 4
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 0 for a python file
func TestFindPrefixIndexSingleLineInPython(t *testing.T) {
	line := fmt.Sprintf("%s single line comment in python", single_line_prefix_python)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return -1 when no single line comment prefix exists
func TestFindPrefixIndexSingleLineNotFoundInPython(t *testing.T) {
	line := "n = 5" // no prefix added
	fields := strings.Fields(line)
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := -1
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// since we didn't find the prefix, the stack should be empty
	require.Len(t, notation.Stack.Items, 0)
}

// should return the prefix index as 3 for a python file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexSingleLineAfterSourceCodeInPython(t *testing.T) {
	line := fmt.Sprintf("n = 1 %s single line comment in python", single_line_prefix_python)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := 3
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 0 for a go file
func TestFindPrefixIndexSingleLineInGo(t *testing.T) {
	line := fmt.Sprintf("%s single line comment in go", single_line_prefix_go)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_go, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 3 for a go file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexSingleLineAfterSourceCodeInGo(t *testing.T) {
	line := fmt.Sprintf("n := 5 %s single line comment in go", single_line_prefix_go)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_go, annotation, &bufio.Scanner{})
	expected := 3
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 0 for a c file
func TestFindPrefixIndexMultiLineInC(t *testing.T) {
	line := fmt.Sprintf("%s multi line comment in c", multi_line_prefix_c)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 4 for a c file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexMultiLineAfterSourceCodeInC(t *testing.T) {
	line := fmt.Sprintf("int main() {return 0}; %s multi line comment in c", multi_line_prefix_c)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := 4
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 0 for a python file
func TestFindPrefixIndexMultiLineInPython(t *testing.T) {
	line := fmt.Sprintf("%s multi line comment in python", multi_line_prefix_python)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 0 for a go file
func TestFindPrefixIndexMultiLineInGo(t *testing.T) {
	line := fmt.Sprintf("%s multi line comment in go", multi_line_prefix_go)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_go, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should return the prefix index as 3 for a go file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexMultiLineAfterSourceCodeInGo(t *testing.T) {
	line := fmt.Sprintf("n := 1 %s multi line comment in go", multi_line_prefix_go)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_go, annotation, &bufio.Scanner{})
	expected := 3
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
	// stack should also have one item in it after locating the prefix
	require.Len(t, notation.Stack.Items, 1)
}

// should extract all text after the annotation for a single line comment in c
func TestExtractFromSingleLineCommentInC(t *testing.T) {
	line := fmt.Sprintf(
		"int main() {return 0} %s %s single line comment",
		multi_line_prefix_c,
		annotation,
	)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := "single line comment"
	start := notation.FindPrefixIndex(fields) // find prefix and increment by 1
	actual := notation.ExtractFromSingleLineComment(fields, start+1)
	require.Equal(t, expected, actual)
}

// should return an empty string when an annotation is not found
func TestExtractFromSingleLineCommentEmptyInC(t *testing.T) {
	line := "int main() {return 0}"
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	expected := ""
	start := notation.FindPrefixIndex(fields) // find prefix and increment by 1
	actual := notation.ExtractFromSingleLineComment(fields, start+1)
	require.Equal(t, expected, actual)
}

// should extract all text after the annotation for a single line comment in python
func TestExtractFromSingleLineCommentInPython(t *testing.T) {
	line := fmt.Sprintf("n = 5 %s %s single line comment", single_line_prefix_python, annotation)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := "single line comment"
	start := notation.FindPrefixIndex(fields) // find prefix and increment by 1
	actual := notation.ExtractFromSingleLineComment(fields, start+1)
	require.Equal(t, expected, actual)
}

// should return an empty string when an annotation is not found
func TestExtractFromSingleLineCommentEmptyInPython(t *testing.T) {
	line := "n = 5"
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := ""
	start := notation.FindPrefixIndex(fields) // find prefix and increment by 1
	actual := notation.ExtractFromSingleLineComment(fields, start+1)
	require.Equal(t, expected, actual)
}
