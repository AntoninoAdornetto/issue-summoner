package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

const (
	test_annotation = "@TEST_TODO"
)

func TestNewIssueManagerNotReporting(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, false)
	require.NoError(t, err)
	require.NotNil(t, manager)
}

func TestNewIssueManagerReporting(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, true)
	require.NoError(t, err)
	require.NotNil(t, manager)
}
