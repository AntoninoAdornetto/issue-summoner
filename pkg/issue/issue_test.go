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
