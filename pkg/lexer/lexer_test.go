package lexer_test

import (
	"io"
	"os"
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

func TestAnalyzeTokensCSrcCode(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected []lexer.Token
	}{
		{
			name: "should return 1 token (EOF) when there are no comments present in a c source code file",
			path: "../../testdata/fixtures/c/no-comments.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_EOF,
					Lexeme: []byte{0},
					Line:   7,
					Start:  112,
					End:    112,
				},
			},
		},
		{
			name: "should return the correct tokens for c source code",
			path: "../../testdata/fixtures/c/mix.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   5,
					Start:  62,
					End:    63,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: testAnnotation,
					Line:   5,
					Start:  65,
					End:    80,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("inline"),
					Line:   5,
					Start:  82,
					End:    87,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   5,
					Start:  89,
					End:    95,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("#1"),
					Line:   5,
					Start:  97,
					End:    98,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Line:   5,
					Start:  100,
					End:    101,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   6,
					Start:  121,
					End:    122,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: testAnnotation,
					Line:   6,
					Start:  124,
					End:    139,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("inline"),
					Line:   6,
					Start:  141,
					End:    146,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   6,
					Start:  148,
					End:    154,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("#2"),
					Line:   6,
					Start:  156,
					End:    157,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/;"),
					Line:   6,
					Start:  159,
					End:    161,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   10,
					Start:  204,
					End:    205,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: testAnnotation,
					Line:   10,
					Start:  207,
					End:    222,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("decode"),
					Line:   10,
					Start:  224,
					End:    229,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("the"),
					Line:   10,
					Start:  231,
					End:    233,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("message"),
					Line:   10,
					Start:  235,
					End:    241,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("and"),
					Line:   10,
					Start:  243,
					End:    245,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("clean"),
					Line:   10,
					Start:  247,
					End:    251,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("up"),
					Line:   10,
					Start:  253,
					End:    254,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("after"),
					Line:   10,
					Start:  256,
					End:    260,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("yourself!"),
					Line:   10,
					Start:  262,
					End:    270,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   11,
					Start:  272,
					End:    272,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   14,
					Start:  287,
					End:    288,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: testAnnotation,
					Line:   15,
					Start:  293,
					End:    308,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("drop"),
					Line:   15,
					Start:  310,
					End:    313,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("a"),
					Line:   15,
					Start:  315,
					End:    315,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("star"),
					Line:   15,
					Start:  317,
					End:    320,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("if"),
					Line:   15,
					Start:  322,
					End:    323,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("you"),
					Line:   15,
					Start:  325,
					End:    327,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("know"),
					Line:   15,
					Start:  329,
					End:    332,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("about"),
					Line:   15,
					Start:  334,
					End:    338,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("this"),
					Line:   15,
					Start:  340,
					End:    343,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("code"),
					Line:   15,
					Start:  345,
					End:    348,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("wars"),
					Line:   15,
					Start:  350,
					End:    353,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("challenge"),
					Line:   15,
					Start:  355,
					End:    363,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Digital"),
					Line:   16,
					Start:  368,
					End:    374,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Cypher"),
					Line:   16,
					Start:  376,
					End:    381,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("assigns"),
					Line:   16,
					Start:  383,
					End:    389,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("to"),
					Line:   16,
					Start:  391,
					End:    392,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("each"),
					Line:   16,
					Start:  394,
					End:    397,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("letter"),
					Line:   16,
					Start:  399,
					End:    404,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("of"),
					Line:   16,
					Start:  406,
					End:    407,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   16,
					Start:  409,
					End:    411,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("alphabet"),
					Line:   16,
					Start:  413,
					End:    420,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("unique"),
					Line:   16,
					Start:  422,
					End:    427,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("number."),
					Line:   16,
					Start:  429,
					End:    435,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Instead"),
					Line:   17,
					Start:  440,
					End:    446,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("of"),
					Line:   17,
					Start:  448,
					End:    449,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("letters"),
					Line:   17,
					Start:  451,
					End:    457,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("in"),
					Line:   17,
					Start:  459,
					End:    460,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("encrypted"),
					Line:   17,
					Start:  462,
					End:    470,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("word"),
					Line:   17,
					Start:  472,
					End:    475,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("we"),
					Line:   17,
					Start:  477,
					End:    478,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("write"),
					Line:   17,
					Start:  480,
					End:    484,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   17,
					Start:  486,
					End:    488,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("corresponding"),
					Line:   17,
					Start:  490,
					End:    502,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("number"),
					Line:   17,
					Start:  504,
					End:    509,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Then"),
					Line:   18,
					Start:  514,
					End:    517,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("we"),
					Line:   18,
					Start:  519,
					End:    520,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("add"),
					Line:   18,
					Start:  522,
					End:    524,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("to"),
					Line:   18,
					Start:  526,
					End:    527,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("each"),
					Line:   18,
					Start:  529,
					End:    532,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("obtained"),
					Line:   18,
					Start:  534,
					End:    541,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("digit"),
					Line:   18,
					Start:  543,
					End:    547,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("consecutive"),
					Line:   18,
					Start:  549,
					End:    559,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("digits"),
					Line:   18,
					Start:  561,
					End:    566,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("from"),
					Line:   18,
					Start:  568,
					End:    571,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   18,
					Start:  573,
					End:    575,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("key"),
					Line:   18,
					Start:  577,
					End:    579,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Line:   19,
					Start:  584,
					End:    585,
				},
				{
					Type:   lexer.TOKEN_EOF,
					Lexeme: []byte{0},
					Line:   33,
					Start:  930,
					End:    930,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := getSrcCode(t, tc.path)
			base := lexer.NewLexer(testAnnotation, src, tc.path)

			target, err := lexer.NewTargetLexer(base)
			require.NoError(t, err)

			tokens, err := base.AnalyzeTokens(target)
			require.NoError(t, err)
			require.Equal(t, tc.expected, tokens)

			for _, token := range tokens {
				if token.Type&lexer.TOKEN_SINGLE_LINE_COMMENT_END != 0 {
					continue
				} else if token.Type&lexer.TOKEN_EOF != 0 {
					continue
				} else {
					require.Equal(t, token.Lexeme, base.Src[token.Start:token.End+1])
				}
			}
		})
	}
}

func getSrcCode(t *testing.T, path string) []byte {
	f, err := os.Open(path)
	require.NoError(t, err)

	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)

	return data
}
