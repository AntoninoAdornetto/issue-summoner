package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestCommentSymbols_C(t *testing.T) {
	m := issue.GetCommentSymbols(".c")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_GO(t *testing.T) {
	m := issue.GetCommentSymbols(".go")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_JS(t *testing.T) {
	m := issue.GetCommentSymbols(".js")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_TS(t *testing.T) {
	m := issue.GetCommentSymbols(".ts")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_TSX(t *testing.T) {
	m := issue.GetCommentSymbols(".tsx")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_JSX(t *testing.T) {
	m := issue.GetCommentSymbols(".jsx")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_PYTHON(t *testing.T) {
	m := issue.GetCommentSymbols(".py")
	expectedSingleLinePrefix := []string{"#"}
	expectedMultiLineStartPrefix := []string{"\"\"\"", "'''"}
	expectedMultiLineEndPrefix := []string{"\"\"\"", "'''"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_MARKDOWN(t *testing.T) {
	m := issue.GetCommentSymbols(".md")
	expectedMultiLineStartPrefix := []string{"<!--"}
	expectedMultiLineEndPrefix := []string{"-->"}
	require.Empty(t, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}
