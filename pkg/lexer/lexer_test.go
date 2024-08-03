package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

var (
	testAnnotation = []byte("@TEST_ANNOTATION")
)

func TestNewLexer(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected *lexer.Lexer
	}{
		{
			name: "should create a new base lexer using c source code",
			path: "../../testdata/fixtures/c/no-comments.c",
			expected: &lexer.Lexer{
				FilePath:   "../../testdata/fixtures/c/no-comments.c",
				FileName:   "no-comments.c",
				Tokens:     make([]lexer.Token, 0),
				Start:      0,
				Current:    0,
				Line:       1,
				Annotation: testAnnotation,
			},
		},
		{
			name: "should create a new base lexer using go source code",
			path: "../../testdata/fixtures/go/no-comments.go",
			expected: &lexer.Lexer{
				FilePath:   "../../testdata/fixtures/go/no-comments.go",
				FileName:   "no-comments.go",
				Tokens:     make([]lexer.Token, 0),
				Start:      0,
				Current:    0,
				Line:       1,
				Annotation: testAnnotation,
			},
		},
		{
			name: "should create a new base lexer using js source code",
			path: "../../testdata/fixtures/js/no-comments.js",
			expected: &lexer.Lexer{
				FilePath:   "../../testdata/fixtures/js/no-comments.js",
				FileName:   "no-comments.js",
				Tokens:     make([]lexer.Token, 0),
				Start:      0,
				Current:    0,
				Line:       1,
				Annotation: testAnnotation,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := getSrcCode(t, tc.path)
			tc.expected.Src = src
			actual := lexer.NewLexer(testAnnotation, src, tc.path)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestNewTargetLexer(t *testing.T) {
	testCases := []struct {
		name     string
		base     *lexer.Lexer
		expected lexer.LexicalTokenizer
	}{
		{
			name: "Should create a c-lexer (target lexer) when provided c source code",
			base: lexer.NewLexer(
				testAnnotation,
				getSrcCode(t, "../../testdata/fixtures/c/no-comments.c"),
				"../../testdata/fixtures/c/no-comments.c",
			),
			expected: &lexer.Clexer{},
		},
		{
			name: "Should create a c-lexer (target lexer) when provided go source code",
			base: lexer.NewLexer(
				testAnnotation,
				getSrcCode(t, "../../testdata/fixtures/go/no-comments.go"),
				"../../testdata/fixtures/go/no-comments.go",
			),
			expected: &lexer.Clexer{},
		},
		{
			name: "Should create a c-lexer (target lexer) when provided js source code",
			base: lexer.NewLexer(
				testAnnotation,
				getSrcCode(t, "../../testdata/fixtures/js/no-comments.js"),
				"../../testdata/fixtures/js/no-comments.js",
			),
			expected: &lexer.Clexer{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := lexer.NewTargetLexer(tc.base)
			require.NoError(t, err)
			require.IsType(t, tc.expected, actual)
		})
	}
}


// should return the lexing manager for c file extension
func TestNewLexingManagerC(t *testing.T) {
	lm, err := lexer.NewLexingManager(".c")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for js file extension
func TestNewLexingManagerJS(t *testing.T) {
	lm, err := lexer.NewLexingManager(".js")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for ts file extension
func TestNewLexingManagerTS(t *testing.T) {
	lm, err := lexer.NewLexingManager(".ts")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for cpp file extension
func TestNewLexingManagerCPP(t *testing.T) {
	lm, err := lexer.NewLexingManager(".cpp")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for java file extension
func TestNewLexingManagerJava(t *testing.T) {
	lm, err := lexer.NewLexingManager(".java")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for go file extension
func TestNewLexingManagerGo(t *testing.T) {
	lm, err := lexer.NewLexingManager(".go")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return an error when an unsupported file extension is provided
func TestNewLexingManagerUnsupported(t *testing.T) {
	lm, err := lexer.NewLexingManager(".unsupported")
	require.Error(t, err)
	require.Nil(t, lm)
}
