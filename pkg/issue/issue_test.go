package issue_test

import (
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
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_SingleLineGo(t *testing.T) {
	line := "// single line comment"
	expected := issue.LINE_TYPE_SINGLE
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_MultiLineStartGo(t *testing.T) {
	line := "/* Start multi line comment"
	expected := issue.LINE_TYPE_MULTI
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_MultiLineEndGo(t *testing.T) {
	line := "*/"
	expected := issue.LINE_TYPE_MULTI
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".go"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_SrcCodePy(t *testing.T) {
	line := "def sum(a: int, b: int):"
	expected := issue.LINE_TYPE_SRC_CODE
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_SingleLinePy(t *testing.T) {
	line := "# single line comment"
	expected := issue.LINE_TYPE_SINGLE
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_MultiLineStartPy(t *testing.T) {
	line := "''' Start multi line comment"
	expected := issue.LINE_TYPE_MULTI
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
}

func TestEvalSourceLine_MultiLineEndPy(t *testing.T) {
	line := "End multi line comment'''"
	expected := issue.LINE_TYPE_MULTI
	actual := issue.EvalSourceLine(line, issue.GetCommentSymbols(".py"))
	require.Equal(t, expected, actual)
}
