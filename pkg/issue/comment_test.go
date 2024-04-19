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
}

// should return the prefix index as 0 for a python file
func TestFindPrefixIndexSingleLineInPython(t *testing.T) {
	line := fmt.Sprintf("%s single line comment in python", single_line_prefix_python)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should return the prefix index as 3 for a python file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexSingleLineAfterSourceCodeInPython(t *testing.T) {
	line := fmt.Sprintf("dict = 1 %s single line comment in python", single_line_prefix_python)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_py, annotation, &bufio.Scanner{})
	expected := 3
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should return the prefix index as 0 for a go file
func TestFindPrefixIndexSingleLineInGo(t *testing.T) {
	line := fmt.Sprintf("%s single line comment in go", single_line_prefix_go)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 0
	notation := issue.NewCommentNotation(file_ext_go, annotation, &bufio.Scanner{})
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should return the prefix index as 3 for a go file where the
// prefix notation is not denoted until after the source code
func TestFindPrefixIndexSingleLineAfterSourceCodeInGo(t *testing.T) {
	line := fmt.Sprintf("dict = 1 %s single line comment in go", single_line_prefix_go)
	fields := strings.Fields(line) // after splitting, the prefix is located at index 4
	notation := issue.NewCommentNotation(file_ext_go, annotation, &bufio.Scanner{})
	expected := 3
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}
