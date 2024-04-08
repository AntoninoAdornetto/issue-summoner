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
