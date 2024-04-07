package issue_test

import (
	"fmt"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestGetIssueManager_Pending(t *testing.T) {
	issueManager, err := issue.GetIssueManager(issue.PENDING_ISSUE, "@TEST_TODO")
	require.NoError(t, err)
	require.IsType(t, &issue.PendingIssue{}, issueManager)
}

func TestGetIssueManager_Processed(t *testing.T) {
	issueManager, err := issue.GetIssueManager(issue.PROCESSED_ISSUE, "@TEST_TODO")
	require.NoError(t, err)
	require.IsType(t, &issue.ProcessedIssue{}, issueManager)
}

func TestGetIssueManager_UnknownIssueType(t *testing.T) {
	issueManager, err := issue.GetIssueManager("unknown", "@TEST_TODO")
	require.Error(t, err)
	require.Nil(t, issueManager)
}

func TestEvalSourceLine_SrcCodeGo(t *testing.T) {
	line := "func main(){ fmt.Printf('Hello World\n')}"
	expected := issue.LINE_TYPE_SRC_CODE
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
	require.Equal(t, "", prefix)
}

func TestEvalSourceLine_SingleLineGo(t *testing.T) {
	line := "// single line comment"
	expected := issue.LINE_TYPE_SINGLE
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
	require.Equal(t, "//", prefix)
}

func TestEvalSourceLine_MultiLineStartGo(t *testing.T) {
	line := "/* Start multi line comment"
	expected := issue.LINE_TYPE_MULTI_START
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
	require.Equal(t, "/*", prefix)
}

func TestEvalSourceLine_MultiLineEndGo(t *testing.T) {
	line := "*/"
	expected := issue.LINE_TYPE_MULTI_END
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
	require.Equal(t, "*/", prefix)
}

func TestEvalSourceLine_SrcCodePy(t *testing.T) {
	line := "def sum(a: int, b: int):"
	expected := issue.LINE_TYPE_SRC_CODE
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
	require.Equal(t, "", prefix)
}

func TestEvalSourceLine_SingleLinePy(t *testing.T) {
	line := "# single line comment"
	expected := issue.LINE_TYPE_SINGLE
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
	require.Equal(t, "#", prefix)
}

func TestEvalSourceLine_MultiLineStartPy(t *testing.T) {
	line := "''' Start multi line comment"
	expected := issue.LINE_TYPE_MULTI_START
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
	require.Equal(t, "'''", prefix)
}

func TestEvalSourceLine_MultiLineEndPy(t *testing.T) {
	line := "End multi line comment'''"
	expected := issue.LINE_TYPE_MULTI_END
	actual, prefix := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
	require.Equal(t, "'''", prefix)
}

// should remove the single line comment prefix, and the annotation (@TEST_TODO) then return all text after the annotation
func TestParseSingleLineCommentWithAnnotationGo(t *testing.T) {
	prefix := "//"
	annotation := "@TEST_TODO"
	expected := "This is a single line comment with an annotation prepended to it"
	line := fmt.Sprintf("%s %s %s", prefix, annotation, expected)
	actual, isAnnotated := issue.ParseSingleLineComment(line, annotation, prefix)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
}

// should remove the single line comment prefix, and then return all text after the prefix
func TestParseSingleLineCommentWithOutAnnotationGo(t *testing.T) {
	prefix := "//"
	annotation := "@TEST_TODO"
	expected := "This is a single line comment without an annotation"
	line := fmt.Sprintf("%s %s", prefix, expected)
	actual, isAnnotated := issue.ParseSingleLineComment(line, annotation, prefix)
	require.Equal(t, expected, actual)
	require.False(t, isAnnotated)
}

// should remove the single line comment prefix, and the annotation (@TEST_TODO) then return all text after the annotation
func TestParseSingleLineCommentWithAnnotationPython(t *testing.T) {
	prefix := "#"
	annotation := "@TEST_TODO"
	expected := "This is a single line comment with an annotation prepended to it"
	line := fmt.Sprintf("%s %s %s", prefix, annotation, expected)
	actual, isAnnotated := issue.ParseSingleLineComment(line, annotation, prefix)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
}

// should remove the single line comment prefix, and then return all text after the prefix
func TestParseSingleLineCommentWithOutAnnotationPython(t *testing.T) {
	prefix := "#"
	annotation := "@TEST_TODO"
	expected := "This is a single line comment without an annotation"
	line := fmt.Sprintf("%s %s", prefix, expected)
	actual, isAnnotated := issue.ParseSingleLineComment(line, annotation, prefix)
	require.Equal(t, expected, actual)
	require.False(t, isAnnotated)
}
