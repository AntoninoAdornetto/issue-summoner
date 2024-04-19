package issue_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

const (
	file_ext_c        = ".c"
	file_ext_go       = ".go"
	file_ext_python   = ".py"
	file_ext_default  = "default"
	file_ext_markdown = ".md"
)

func TestNewCommentNotationC(t *testing.T) {
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	require.Equal(t, notation.SingleLinePrefix, issue.CommentNotations[file_ext_c].SingleLinePrefix)
	require.Equal(t, notation.SingleLineSuffix, issue.CommentNotations[file_ext_c].SingleLineSuffix)
	require.Equal(t, notation.MultiLinePrefix, issue.CommentNotations[file_ext_c].MultiLinePrefix)
	require.Equal(t, notation.MultiLineSuffix, issue.CommentNotations[file_ext_c].MultiLineSuffix)
	require.Equal(t, notation.Annotation, annotation)
}

func TestNewCommentNotationPython(t *testing.T) {
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	require.Equal(
		t,
		notation.SingleLinePrefix,
		issue.CommentNotations[file_ext_python].SingleLinePrefix,
	)
	require.Equal(
		t,
		notation.SingleLineSuffix,
		issue.CommentNotations[file_ext_python].SingleLineSuffix,
	)
	require.Equal(
		t,
		notation.MultiLinePrefix,
		issue.CommentNotations[file_ext_python].MultiLinePrefix,
	)
	require.Equal(
		t,
		notation.MultiLineSuffix,
		issue.CommentNotations[file_ext_python].MultiLineSuffix,
	)
	require.Equal(t, notation.Annotation, annotation)
}

func TestNewCommentNotationMarkdown(t *testing.T) {
	notation := issue.NewCommentNotation(file_ext_markdown, annotation, nil)
	require.Equal(
		t,
		notation.SingleLinePrefix,
		issue.CommentNotations[file_ext_markdown].SingleLinePrefix,
	)
	require.Equal(
		t,
		notation.SingleLineSuffix,
		issue.CommentNotations[file_ext_markdown].SingleLineSuffix,
	)
	require.Equal(t, notation.Annotation, annotation)
}

func TestNewCommentNotationDefault(t *testing.T) {
	notation := issue.NewCommentNotation(file_ext_default, annotation, nil)
	require.Equal(
		t,
		notation.SingleLinePrefix,
		issue.CommentNotations[file_ext_default].SingleLinePrefix,
	)
	require.Equal(t, notation.Annotation, annotation)
}

func TestParseLineEmpty(t *testing.T) {
	lineNumber := uint64(0)
	notation := issue.NewCommentNotation(file_ext_c, annotation, &bufio.Scanner{})
	issue, err := notation.ParseLine(&lineNumber)
	require.NoError(t, err)
	require.Empty(t, issue)
}

// there should be more tests for this function but am facing issues with the scanner.Scan() method
// so I will write full integration tests using the Walk/Scan methods from the pending/processed issue interfaces
func TestParseLineSingleLineInC(t *testing.T) {
	lineNumber := uint64(3)
	srcLines := []byte(`
	#include <stdio.h>

	// @TEST_TODO This is a test comment
	int main() {
		printf("Hello, World!");
		return 0;
	}
	`)
	r := bufio.NewReader(bytes.NewReader(srcLines))
	scanner := bufio.NewScanner(r)

	// there is a function that will be in charge of calling the scanner.Scan() method
	// for this test we will call it manually but in the implementation it will be called by the Walk/Scan method
	// of pending/processed issue interface implementations
	for i := 0; i <= 3; i++ {
		scanner.Scan()
	}

	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      3,
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}

	notation := issue.NewCommentNotation(file_ext_c, annotation, scanner)
	actual, err := notation.ParseLine(&lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a go file
func TestFindPrefixIndexSingleLineCommentForGoFile(t *testing.T) {
	fields := []string{"//", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_go, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a go file when the line does not begin with the comment prefix
func TestFindPrefixIndexSingleLineCommentAfterGoCodeTokens(t *testing.T) {
	fields := []string{
		"package",
		"main",
		"import",
		"fmt",
		"//",
		annotation,
		"This",
		"is",
		"a",
		"test",
		"comment",
	}
	notation := issue.NewCommentNotation(file_ext_go, annotation, nil)
	expected := 4
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a c file
func TestFindPrefixIndexSingleLineCommentForCFile(t *testing.T) {
	fields := []string{"//", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a c file when the line does not begin with the comment prefix
func TestFindPrefixIndexSingleLineCommentAfterCCodeTokens(t *testing.T) {
	fields := []string{
		"int",
		"main",
		"()",
		"{",
		"//",
		annotation,
		"This",
		"is",
		"a",
		"test",
		"comment",
		"}",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := 4
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a python file
func TestFindPrefixIndexSingleLineCommentForPythonFile(t *testing.T) {
	fields := []string{"#", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a python file when the line does not begin with the comment prefix
func TestFindPrefixIndexSingleLineCommentAfterPythonCodeTokens(t *testing.T) {
	fields := []string{"def", "main", "():", "#", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := 3
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the single line comment prefix for a markdown file
func TestFindPrefixIndexSingleLineCommentForMarkdownFile(t *testing.T) {
	fields := []string{"<!--", annotation, "This", "is", "a", "test", "comment", "-->"}
	notation := issue.NewCommentNotation(file_ext_markdown, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the multi line comment prefix for a c file
func TestFindPrefixIndexMultiLineCommentForCFile(t *testing.T) {
	fields := []string{"/*", annotation, "This", "is", "a", "test", "comment", "*/"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the multi line comment prefix for a c file when the line does not begin with the comment prefix
func TestFindPrefixIndexMultiLineCommentAfterCCodeTokens(t *testing.T) {
	fields := []string{
		"int",
		"main",
		"()",
		"{",
		"/*",
		annotation,
		"This",
		"is",
		"a",
		"test",
		"comment",
		"*/",
		"}",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := 4
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the multi line comment prefix for a go file
func TestFindPrefixIndexMultiLineCommentForGoFile(t *testing.T) {
	fields := []string{"/*", annotation, "This", "is", "a", "test", "comment", "*/"}
	notation := issue.NewCommentNotation(file_ext_go, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the multi line comment prefix for a go file when the line does not begin with the comment prefix
func TestFindPrefixIndexMultiLineCommentAfterGoCodeTokens(t *testing.T) {
	fields := []string{
		"package",
		"main",
		"import",
		"fmt",
		"/*",
		annotation,
		"This",
		"is",
		"a",
		"test",
		"comment",
		"*/",
	}
	notation := issue.NewCommentNotation(file_ext_go, annotation, nil)
	expected := 4
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should locate the index of the multi line comment prefix for a python file
func TestFindPrefixIndexMultiLineCommentForPythonFile(t *testing.T) {
	fields := []string{"'''", annotation, "This", "is", "a", "test", "comment", "'''"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := 0
	actual := notation.FindPrefixIndex(fields)
	require.Equal(t, expected, actual)
}

// should extract the comment from a multi line comment in a C file.
// note: the comment is the start of the multi line comment but the end is not present for this test
func TestExtractFromMultiPrefixCommentForCFile(t *testing.T) {
	fields := []string{"/*", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.True(
		t,
		notation.AnnotationIndicator,
	) // this line contained the annotation so this should be true
}

// should extract all text after the new line comment prefix symbol (*)
// note: this is a continuation of the previous test
func TestExtractFromMultiLineNewLineCommentForCFile(t *testing.T) {
	fields := []string{"*", "continue", "comment", "from", "the", "previous", "test"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := "continue comment from the previous test"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extract all the text prior to the suffix
func TestExtractFromMultiLineSuffixCommentForCFile(t *testing.T) {
	fields := []string{
		"this",
		"ends",
		"the",
		"comment",
		"from",
		"the",
		"previous",
		"2",
		"tests",
		"*/",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	// push onto the stack because once a suffix is found, the function will pop the stack. This is to ensure that the stack is not empty
	notation.Stack.Push(notation.MultiLinePrefix)
	expected := "this ends the comment from the previous 2 tests"
	// we increment by 1 because there is no prefix index, it returns -1 which means we want to start at beginging of the line
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return an empty string when the suffix is the only symbol in the line
func TestExtractFromMultiLineSuffixCommentForCFileWithSuffixOnly(t *testing.T) {
	fields := []string{"*/"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	// push onto the stack because once a suffix is found, the function will pop the stack. This is to ensure that the stack is not empty
	notation.Stack.Push(notation.MultiLinePrefix)
	expected := ""
	// we increment by 1 because there is no prefix index, it returns -1 which means we want to start at beginging of the line
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extract the comment from a multi line comment in a python file.
// note: the comment is the start of the multi line comment but the end is not present for this test
func TestExtractFromMultiPrefixCommentForPythonFile(t *testing.T) {
	fields := []string{"'''", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extract all the text since the multi line comment has not been closed yet
// note: this is a continuation of the previous test
func TestExtractFromMultiLineNewLineCommentForPythonFile(t *testing.T) {
	fields := []string{"continue", "comment", "from", "the", "previous", "test"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := "continue comment from the previous test"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestExtractFromMultiLineSuffixCommentForPythonFile(t *testing.T) {
	fields := []string{
		"this",
		"ends",
		"the",
		"comment",
		"from",
		"the",
		"previous",
		"2",
		"tests",
		"'''",
	}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	// push onto the stack because once a suffix is found, the function will pop the stack. This is to ensure that the stack is not empty
	notation.Stack.Push(notation.MultiLinePrefix)
	expected := "this ends the comment from the previous 2 tests"
	// hard coding to -1 because some languages denote both the prefix and suffix with the same symbol
	// we won't ever have to worry about this in the implementation because the prefix index will always be 0 once the prefix is found
	// this is just to test the functionality of the function
	prefixIndex := -1
	actual, err := notation.ExtractFromMulti(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extract all text after the prefix and annotation
func TestExtractFromSingleLineForCFile(t *testing.T) {
	fields := []string{"//", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromSingle(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extract all text after source code tokens and the prefix and annotation
func TestExtractFromSingleLineForCFileAfterTokens(t *testing.T) {
	fields := []string{
		"int",
		"main",
		"()",
		"{",
		"}",
		"//",
		annotation,
		"This",
		"is",
		"a",
		"test",
		"comment",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromSingle(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extact all text after the prefix and annotation for a python file
func TestExtractFromSingleLineForPythonFile(t *testing.T) {
	fields := []string{"#", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromSingle(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestExtractFromSingleLineForPythonFileAfterTokens(t *testing.T) {
	fields := []string{"def", "main", "():", "#", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromSingle(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should extract all text between the prefix and suffix for a markdown file
func TestExtractFromSingleLineForMarkdownFile(t *testing.T) {
	fields := []string{"<!--", annotation, "This", "is", "a", "test", "comment", "-->"}
	notation := issue.NewCommentNotation(file_ext_markdown, annotation, nil)
	expected := "This is a test comment"
	prefixIndex := notation.FindPrefixIndex(fields)
	actual, err := notation.ExtractFromSingle(fields, prefixIndex+1)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should build an issue based off a single line comment for a c file
func TestBuildSingleIssueFromSingleLineCommentForCFile(t *testing.T) {
	fields := []string{"//", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      3,
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildSingle(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestBuildSingleIssueFromSingleLineCommentForCFileAfterTokens(t *testing.T) {
	fields := []string{
		"int",
		"main",
		"()",
		"{",
		"}",
		"//",
		annotation,
		"This",
		"is",
		"a",
		"test",
		"comment",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      3,
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildSingle(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should build an issue based off a single line comment for a python file
func TestBuildSingleIssueFromSingleLineCommentForPythonFile(t *testing.T) {
	fields := []string{"#", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      3,
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildSingle(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should build an issue based off a single line comment for a python file
func TestBuildSingleIssueFromSingleLineCommentForPythonFileAfterTokens(t *testing.T) {
	fields := []string{"def", "main", "():", "#", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_python, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      3,
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildSingle(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestBuildSingleIssueFromSingleLineCommentForMarkdownFile(t *testing.T) {
	fields := []string{"<!--", annotation, "This", "is", "a", "test", "comment", "-->"}
	notation := issue.NewCommentNotation(file_ext_markdown, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      3,
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildSingle(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should return an error when a prefix is not found and the stack is empty
func TestBuildSingleIssueEmptyStack(t *testing.T) {
	fields := []string{annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	lineNumber := uint64(3)
	actual := issue.Issue{}
	err := notation.BuildSingle(fields, &actual, 1, &lineNumber)
	require.Errorf(t, err, "error: notation stack underflow")
}

// should start building out the issue for a multi line comment.
// note: the next two tests will be a continuation of this test
func TestBuildFromMultiLinePrefixWithAnnotationCommentForCFile(t *testing.T) {
	fields := []string{"/*", annotation, "This", "is", "a", "test", "comment"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "",
		StartLineNumber:      0, // this is at 0 because the parseLine function will actually set the start line number
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildMulti(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.True(t, notation.AnnotationIndicator)
}

// should build off the issue from the previous test and continue to build the issue
func TestBuildFromMultiNewLineCommentForCFile(t *testing.T) {
	fields := []string{"*", "continue", "comment", "from", "the", "previous", "test"}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "continue comment from the previous test",
		StartLineNumber:      0, // this is at 0 because the parseLine function will actually set the start line number
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{
		Title:                "This is a test comment",
		StartLineNumber:      0,
		AnnotationLineNumber: 3,
	}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildMulti(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// should build off the issue from the previous test and finalize the issue now that the suffix has been found
func TestBuildFromMultiLineSuffixForCFile(t *testing.T) {
	fields := []string{
		"this",
		"ends",
		"the",
		"comment",
		"from",
		"the",
		"previous",
		"2",
		"tests",
		"*/",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	notation.Stack.Push(
		notation.MultiLinePrefix,
	) // have to push because once a suffix is found, we pop the stack
	lineNumber := uint64(3)
	expected := issue.Issue{
		Title:                "This is a test comment",
		Description:          "continue comment from the previous test this ends the comment from the previous 2 tests",
		StartLineNumber:      0, // this is at 0 because the parseLine function will actually set the start line number
		EndLineNumber:        3,
		AnnotationLineNumber: 3,
	}
	actual := issue.Issue{
		Title:                "This is a test comment",
		Description:          "continue comment from the previous test",
		StartLineNumber:      0,
		AnnotationLineNumber: 3,
	}
	prefixIndex := notation.FindPrefixIndex(fields)
	err := notation.BuildMulti(fields, &actual, prefixIndex+1, &lineNumber)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestBuildFromMultiLineEmptyStackError(t *testing.T) {
	fields := []string{
		"this",
		"ends",
		"the",
		"comment",
		"from",
		"the",
		"previous",
		"2",
		"tests",
		"*/",
	}
	notation := issue.NewCommentNotation(file_ext_c, annotation, nil)
	lineNumber := uint64(3)
	actual := issue.Issue{}
	err := notation.BuildMulti(fields, &actual, 1, &lineNumber)
	require.Errorf(t, err, "error: notation stack underflow")
}

// should push items into the stack and increase the Top property
func TestNotationStackPush(t *testing.T) {
	n := issue.NewCommentNotation(".c", annotation, nil).Stack
	items := []string{"//", "#", "/*", "///"}

	// The stack's Top property starts at -1
	for i, item := range items {
		n.Push(item)
		expected := i // 0, 1, 2, 3
		actual := n.Top
		require.Equal(t, expected, actual)
	}
	require.Equal(t, items, n.Items)
}

// should return the top item and decrement the Top property
func TestNotationStackPop(t *testing.T) {
	n := issue.NewCommentNotation(".c", annotation, nil).Stack
	items := []string{"//", "#", "/*", "///"}

	// The stack's Top property starts at -1
	for i, item := range items {
		n.Push(item)
		expected := i // 0, 1, 2, 3
		actual := n.Top
		require.Equal(t, expected, actual)
	}

	for i := len(items) - 1; i <= 0; i-- {
		item, err := n.Pop()
		expected := i // 3, 2, 1, 0
		actual := n.Top
		require.NoError(t, err)
		require.Equal(t, items[i], item)
		require.Equal(t, expected, actual)
	}
}

// should return if the stack is empty or not
func TestNotationStackIsEmpty(t *testing.T) {
	n := issue.NewCommentNotation(".c", annotation, nil).Stack

	// Stack is initially empty
	require.True(t, n.IsEmpty())

	n.Push("//")
	require.False(t, n.IsEmpty())

	_, _ = n.Pop()
	require.True(t, n.IsEmpty())
}

// should return the item at the top
func TestNotationStackGetTopItem(t *testing.T) {
	n := issue.NewCommentNotation(".c", annotation, nil).Stack

	// Stack is initially empty
	_, err := n.Peek()
	require.Error(t, err)

	n.Push("//")
	item, err := n.Peek()
	require.NoError(t, err)
	require.Equal(t, "//", item)

	n.Push("#")
	item, err = n.Peek()
	require.NoError(t, err)
	require.Equal(t, "#", item)

	_, _ = n.Pop()
	item, err = n.Peek()
	require.NoError(t, err)
	require.Equal(t, "//", item)
}
