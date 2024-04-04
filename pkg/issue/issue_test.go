package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestGetIssueManager_Pending(t *testing.T) {
	issueManager, err := issue.GetIssueManager(issue.PENDING_ISSUE)
	require.NoError(t, err)
	require.IsType(t, &issue.PendingIssue{}, issueManager)
}

func TestGetIssueManager_Processed(t *testing.T) {
	issueManager, err := issue.GetIssueManager(issue.PROCESSED_ISSUE)
	require.NoError(t, err)
	require.IsType(t, &issue.ProcessedIssue{}, issueManager)
}

func TestGetIssueManager_UnknownIssueType(t *testing.T) {
	issueManager, err := issue.GetIssueManager("unknownIssueType")
	require.Error(t, err)
	require.Nil(t, issueManager)
}

func TestSkip_SingleLineTrueGo(t *testing.T) {
	line := "func main(){ fmt.Printf('Hello World\n')}"
	require.True(t, issue.Skip(line, issue.GetCommentSymbols(".go")))
}

func TestSkip_SingleLineFalseGo(t *testing.T) {
	line := "// single line comment"
	require.False(t, issue.Skip(line, issue.GetCommentSymbols(".go")))
}

func TestSkip_SingleLineTruePy(t *testing.T) {
	line := "def sum(a: int, b: int):"
	require.True(t, issue.Skip(line, issue.GetCommentSymbols(".py")))
}

func TestSkip_SingleLineFalsePy(t *testing.T) {
	line := "# single line comment"
	require.False(t, issue.Skip(line, issue.GetCommentSymbols(".py")))
}

func TestSkip_MultiLineStartFalseGo(t *testing.T) {
	line := "/* begin multi line comment"
	require.False(t, issue.Skip(line, issue.GetCommentSymbols(".go")))
}

func TestSkip_MultiLineEndFalseGo(t *testing.T) {
	line := "End Multi line comment*/"
	require.False(t, issue.Skip(line, issue.GetCommentSymbols(".go")))
}

func TestSkip_MultiLineStartFalsePy(t *testing.T) {
	line := "''' Start Multi line"
	require.False(t, issue.Skip(line, issue.GetCommentSymbols(".py")))
}

func TestSkip_MultiLineEndFalsePy(t *testing.T) {
	line := "End Multi line comment'''"
	require.False(t, issue.Skip(line, issue.GetCommentSymbols(".py")))
}
