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
