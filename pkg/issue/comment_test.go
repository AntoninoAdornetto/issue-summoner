package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestCommentSymbols_C(t *testing.T) {
	m := issue.GetCommentSymbols(".c")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_GO(t *testing.T) {
	m := issue.GetCommentSymbols(".go")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_JS(t *testing.T) {
	m := issue.GetCommentSymbols(".js")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_TS(t *testing.T) {
	m := issue.GetCommentSymbols(".ts")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_TSX(t *testing.T) {
	m := issue.GetCommentSymbols(".tsx")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_JSX(t *testing.T) {
	m := issue.GetCommentSymbols(".jsx")
	expectedSingleLineSymbols := []string{"//"}
	expectedMultiLineStartSymbols := []string{"/*"}
	expectedMultiLineEndSymbols := []string{"*/"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_PYTHON(t *testing.T) {
	m := issue.GetCommentSymbols(".py")
	expectedSingleLineSymbols := []string{"#"}
	expectedMultiLineStartSymbols := []string{"\"\"\"", "'''"}
	expectedMultiLineEndSymbols := []string{"\"\"\"", "'''"}
	require.Equal(t, expectedSingleLineSymbols, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}

func TestCommentSymbols_MARKDOWN(t *testing.T) {
	m := issue.GetCommentSymbols(".md")
	expectedMultiLineStartSymbols := []string{"<!--"}
	expectedMultiLineEndSymbols := []string{"-->"}
	require.Empty(t, m.SingleLineSymbols)
	require.Equal(t, expectedMultiLineStartSymbols, m.MultiLineStartSymbols)
	require.Equal(t, expectedMultiLineEndSymbols, m.MultiLineEndSymbols)
}
