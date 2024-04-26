package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

const (
	annotation = "@TEST_TODO"
)

func TestNewIssueManagerUnsupported(t *testing.T) {
	im, err := issue.NewIssueManager("unsupported", annotation)
	require.Errorf(t, err, "Unsupported issue type. Use 'pending' or 'processed'")
	require.Nil(t, im)
}

func TestNewIssueManagerPending(t *testing.T) {
	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.IsType(t, &issue.PendingIssue{}, im)
}

func TestNewIssueManagerProcessed(t *testing.T) {
	im, err := issue.NewIssueManager(issue.PROCESSED_ISSUE, annotation)
	require.NoError(t, err)
	require.IsType(t, &issue.ProcessedIssue{}, im)
}
