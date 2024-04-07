package issue_test

import (
	"fmt"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

func TestCommentSymbols_C(t *testing.T) {
	m := issue.CommentPrefixes(".c")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_GO(t *testing.T) {
	m := issue.CommentPrefixes(".go")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_JS(t *testing.T) {
	m := issue.CommentPrefixes(".js")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_TS(t *testing.T) {
	m := issue.CommentPrefixes(".ts")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_TSX(t *testing.T) {
	m := issue.CommentPrefixes(".tsx")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_JSX(t *testing.T) {
	m := issue.CommentPrefixes(".jsx")
	expectedSingleLinePrefix := []string{"//"}
	expectedMultiLineStartPrefix := []string{"/*"}
	expectedMultiLineEndPrefix := []string{"*/"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_PYTHON(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	expectedSingleLinePrefix := []string{"#"}
	expectedMultiLineStartPrefix := []string{"\"\"\"", "'''"}
	expectedMultiLineEndPrefix := []string{"\"\"\"", "'''"}
	require.Equal(t, expectedSingleLinePrefix, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestCommentSymbols_MARKDOWN(t *testing.T) {
	m := issue.CommentPrefixes(".md")
	expectedMultiLineStartPrefix := []string{"<!--"}
	expectedMultiLineEndPrefix := []string{"-->"}
	require.Empty(t, m.SingleLinePrefix)
	require.Equal(t, expectedMultiLineStartPrefix, m.MultiLineStartPrefix)
	require.Equal(t, expectedMultiLineEndPrefix, m.MultiLineEndPrefix)
}

func TestParseLineCommentSingleLinePrefixC(t *testing.T) {
	m := issue.CommentPrefixes(".c")
	prefix := "//"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s Single line comment text", prefix, annotation)
	expected := "Single line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentSingleLinePrefixGo(t *testing.T) {
	m := issue.CommentPrefixes(".go")
	prefix := "//"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s Single line comment text", prefix, annotation)
	expected := "Single line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentSingleLinePrefixPython(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	prefix := "#"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s Single line comment text", prefix, annotation)
	expected := "Single line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentMultiStartLinePrefixC(t *testing.T) {
	m := issue.CommentPrefixes(".c")
	prefix := "/*"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s MultiStart line comment text", prefix, annotation)
	expected := "MultiStart line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentMultiStartLinePrefixGo(t *testing.T) {
	m := issue.CommentPrefixes(".go")
	prefix := "/*"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s MultiStart line comment text", prefix, annotation)
	expected := "MultiStart line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentMultiStartLinePrefixPython(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	prefix := "'''"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s MultiStart line comment text", prefix, annotation)
	expected := "MultiStart line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentMultiStartLinePrefixPythonDoubleQuotes(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	prefix := `"""`
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s %s MultiStart line comment text", prefix, annotation)
	expected := "MultiStart line comment text"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, prefix, m.CurrentPrefix)
}

func TestParseLineCommentMultiEndLinePrefixC(t *testing.T) {
	m := issue.CommentPrefixes(".c")
	suffix := "*/"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("Mutli line end Comment %s", suffix)
	expected := "Mutli line end Comment"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.False(t, isAnnotated)
	require.Equal(t, suffix, m.CurrentPrefix)
}

func TestParseLineCommentMultiEndLinePrefixAnnotationC(t *testing.T) {
	m := issue.CommentPrefixes(".c")
	suffix := "*/"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s Multi line end comment with annotation %s", annotation, suffix)
	expected := "Multi line end comment with annotation"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, suffix, m.CurrentPrefix)
}

func TestParseLineCommentMultiEndLinePrefixPython(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	suffix := "'''"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("Mutli line end Comment %s", suffix)
	expected := "Mutli line end Comment"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.False(t, isAnnotated)
	require.Equal(t, suffix, m.CurrentPrefix)
}

func TestParseLineCommentMultiEndLinePrefixAnnotationPython(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	suffix := "'''"
	annotation := "@TEST_TODO"
	line := fmt.Sprintf("%s Mutli line end Comment %s", annotation, suffix)
	expected := "Mutli line end Comment"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.True(t, isAnnotated)
	require.Equal(t, suffix, m.CurrentPrefix)
}

func TestParseLineNoPrefixOrSuffixC(t *testing.T) {
	m := issue.CommentPrefixes(".c")
	annotation := "@TEST_TODO"
	line := "// no annotation"
	expected := "no annotation"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.False(t, isAnnotated)
}

func TestParseLineNoPrefixOrSuffixPy(t *testing.T) {
	m := issue.CommentPrefixes(".py")
	annotation := "@TEST_TODO"
	line := "# no annotation"
	expected := "no annotation"
	actual, isAnnotated := m.ParseLineComment(line, annotation)
	require.Equal(t, expected, actual)
	require.False(t, isAnnotated)
}
