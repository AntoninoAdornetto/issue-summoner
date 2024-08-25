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
		flags    lexer.U8
	}{
		{
			name:  "should create a new base lexer using c source code",
			path:  "./testdata/c/no-comments.c",
			flags: lexer.FLAG_SCAN,
			expected: &lexer.Lexer{
				FilePath:   "./testdata/c/no-comments.c",
				FileName:   "no-comments.c",
				Tokens:     make([]lexer.Token, 0),
				Start:      0,
				Current:    0,
				Line:       1,
				Annotation: testAnnotation,
			},
		},
		{
			name:  "should create a new base lexer using go source code",
			path:  "./testdata/go/no-comments.go",
			flags: lexer.FLAG_SCAN,
			expected: &lexer.Lexer{
				FilePath:   "./testdata/go/no-comments.go",
				FileName:   "no-comments.go",
				Tokens:     make([]lexer.Token, 0),
				Start:      0,
				Current:    0,
				Line:       1,
				Annotation: testAnnotation,
			},
		},
		{
			name:  "should create a new base lexer using js source code",
			path:  "./testdata/js/no-comments.js",
			flags: lexer.FLAG_SCAN,
			expected: &lexer.Lexer{
				FilePath:   "./testdata/js/no-comments.js",
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
			actual := lexer.NewLexer(testAnnotation, src, tc.path, tc.flags)
			require.Equal(t, tc.expected.Src, actual.Src)
			require.Equal(t, tc.expected.FilePath, actual.FilePath)
			require.Equal(t, tc.expected.FileName, actual.FileName)
			require.Equal(t, tc.expected.Tokens, actual.Tokens)
			require.Equal(t, tc.expected.Start, actual.Start)
			require.Equal(t, tc.expected.Current, actual.Current)
			require.Equal(t, tc.expected.Line, actual.Line)
			require.Equal(t, tc.expected.Annotation, actual.Annotation)
		})
	}
}

func TestNewTargetLexer(t *testing.T) {
	testCases := []struct {
		name     string
		base     *lexer.Lexer
		expected lexer.LexicalTokenizer
		flags    lexer.U8
	}{
		{
			name:  "Should create a c-lexer (target lexer) when provided c source code",
			flags: lexer.FLAG_SCAN,
			base: lexer.NewLexer(
				testAnnotation,
				getSrcCode(t, "./testdata/c/no-comments.c"),
				"./testdata/c/no-comments.c",
				lexer.FLAG_SCAN,
			),
			expected: &lexer.Clexer{},
		},
		{
			name:  "Should create a c-lexer (target lexer) when provided go source code",
			flags: lexer.FLAG_SCAN,
			base: lexer.NewLexer(
				testAnnotation,
				getSrcCode(t, "./testdata/go/no-comments.go"),
				"./testdata/go/no-comments.go",
				lexer.FLAG_SCAN,
			),
			expected: &lexer.Clexer{},
		},
		{
			name:  "Should create a c-lexer (target lexer) when provided js source code",
			flags: lexer.FLAG_SCAN,
			base: lexer.NewLexer(
				testAnnotation,
				getSrcCode(t, "./testdata/js/no-comments.js"),
				"./testdata/js/no-comments.js",
				lexer.FLAG_SCAN,
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
		flags    lexer.U8
	}{
		{
			name:  "should return 1 token (EOF) when there are no comments present in a c source code file",
			path:  "./testdata/c/no-comments.c",
			flags: lexer.FLAG_SCAN,
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
			name:  "should return the correct tokens for c source code",
			path:  "./testdata/c/mix.c",
			flags: lexer.FLAG_SCAN,
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
					Lexeme: []byte("*/"),
					Line:   6,
					Start:  159,
					End:    160,
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
					Lexeme: []byte{'\n'},
					Line:   11,
					Start:  271,
					End:    271,
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
			base := lexer.NewLexer(testAnnotation, src, tc.path, tc.flags)

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

func TestAnalyzeTokensGoSrcCode(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected []lexer.Token
		flags    lexer.U8
	}{
		{
			name:  "should return 1 token (EOF) when there are no comments present in a go source code file",
			path:  "./testdata/go/no-comments.go",
			flags: lexer.FLAG_SCAN,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_EOF,
					Lexeme: []byte{0},
					Line:   23,
					Start:  262,
					End:    262,
				},
			},
		},
		{
			name:  "should return the correct tokens for go source code",
			path:  "./testdata/go/mix.go",
			flags: lexer.FLAG_SCAN,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   9,
					Start:  72,
					End:    73,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Line:   9,
					Start:  75,
					End:    90,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("inline"),
					Line:   9,
					Start:  92,
					End:    97,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   9,
					Start:  99,
					End:    105,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("#1"),
					Line:   9,
					Start:  107,
					End:    108,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Line:   9,
					Start:  110,
					End:    111,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   10,
					Start:  130,
					End:    131,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Line:   10,
					Start:  133,
					End:    148,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("inline"),
					Line:   10,
					Start:  150,
					End:    155,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   10,
					Start:  157,
					End:    163,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("#2"),
					Line:   10,
					Start:  165,
					End:    166,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Line:   10,
					Start:  168,
					End:    169,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   14,
					Start:  192,
					End:    193,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Line:   14,
					Start:  195,
					End:    210,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("decode"),
					Line:   14,
					Start:  212,
					End:    217,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("the"),
					Line:   14,
					Start:  219,
					End:    221,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("message"),
					Line:   14,
					Start:  223,
					End:    229,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("and"),
					Line:   14,
					Start:  231,
					End:    233,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("clean"),
					Line:   14,
					Start:  235,
					End:    239,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("up"),
					Line:   14,
					Start:  241,
					End:    242,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("after"),
					Line:   14,
					Start:  244,
					End:    248,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("yourself!"),
					Line:   14,
					Start:  250,
					End:    258,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{'\n'},
					Line:   15,
					Start:  259,
					End:    259,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   18,
					Start:  273,
					End:    274,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Line:   19,
					Start:  279,
					End:    294,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("drop"),
					Line:   19,
					Start:  296,
					End:    299,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("a"),
					Line:   19,
					Start:  301,
					End:    301,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("star"),
					Line:   19,
					Start:  303,
					End:    306,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("if"),
					Line:   19,
					Start:  308,
					End:    309,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("you"),
					Line:   19,
					Start:  311,
					End:    313,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("know"),
					Line:   19,
					Start:  315,
					End:    318,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("about"),
					Line:   19,
					Start:  320,
					End:    324,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("this"),
					Line:   19,
					Start:  326,
					End:    329,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("code"),
					Line:   19,
					Start:  331,
					End:    334,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("wars"),
					Line:   19,
					Start:  336,
					End:    339,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("challenge"),
					Line:   19,
					Start:  341,
					End:    349,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Digital"),
					Line:   20,
					Start:  354,
					End:    360,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Cypher"),
					Line:   20,
					Start:  362,
					End:    367,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("assigns"),
					Line:   20,
					Start:  369,
					End:    375,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("to"),
					Line:   20,
					Start:  377,
					End:    378,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("each"),
					Line:   20,
					Start:  380,
					End:    383,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("letter"),
					Line:   20,
					Start:  385,
					End:    390,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("of"),
					Line:   20,
					Start:  392,
					End:    393,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   20,
					Start:  395,
					End:    397,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("alphabet"),
					Line:   20,
					Start:  399,
					End:    406,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("unique"),
					Line:   20,
					Start:  408,
					End:    413,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("number."),
					Line:   20,
					Start:  415,
					End:    421,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Instead"),
					Line:   21,
					Start:  426,
					End:    432,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("of"),
					Line:   21,
					Start:  434,
					End:    435,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("letters"),
					Line:   21,
					Start:  437,
					End:    443,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("in"),
					Line:   21,
					Start:  445,
					End:    446,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("encrypted"),
					Line:   21,
					Start:  448,
					End:    456,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("word"),
					Line:   21,
					Start:  458,
					End:    461,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("we"),
					Line:   21,
					Start:  463,
					End:    464,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("write"),
					Line:   21,
					Start:  466,
					End:    470,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   21,
					Start:  472,
					End:    474,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("corresponding"),
					Line:   21,
					Start:  476,
					End:    488,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("number"),
					Line:   21,
					Start:  490,
					End:    495,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Then"),
					Line:   22,
					Start:  500,
					End:    503,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("we"),
					Line:   22,
					Start:  505,
					End:    506,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("add"),
					Line:   22,
					Start:  508,
					End:    510,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("to"),
					Line:   22,
					Start:  512,
					End:    513,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("each"),
					Line:   22,
					Start:  515,
					End:    518,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("obtained"),
					Line:   22,
					Start:  520,
					End:    527,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("digit"),
					Line:   22,
					Start:  529,
					End:    533,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("consecutive"),
					Line:   22,
					Start:  535,
					End:    545,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("digits"),
					Line:   22,
					Start:  547,
					End:    552,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("from"),
					Line:   22,
					Start:  554,
					End:    557,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   22,
					Start:  559,
					End:    561,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("key"),
					Line:   22,
					Start:  563,
					End:    565,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Line:   23,
					Start:  570,
					End:    571,
				},
				{
					Type:   lexer.TOKEN_EOF,
					Lexeme: []byte{0},
					Line:   38,
					Start:  889,
					End:    889,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := getSrcCode(t, tc.path)
			base := lexer.NewLexer(testAnnotation, src, tc.path, tc.flags)

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

func TestAnalyzeTokensJsSrcCode(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected []lexer.Token
		flags    lexer.U8
	}{
		{
			name:  "should do something",
			flags: lexer.FLAG_SCAN,
			path:  "./testdata/js/mix.js",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   1,
					Start:  26,
					End:    27,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: testAnnotation,
					Line:   1,
					Start:  29,
					End:    44,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("fix"),
					Line:   1,
					Start:  46,
					End:    48,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("bug"),
					Line:   1,
					Start:  50,
					End:    52,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("in"),
					Line:   1,
					Start:  54,
					End:    55,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("v8"),
					Line:   1,
					Start:  57,
					End:    58,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{'\n'},
					Line:   2,
					Start:  59,
					End:    59,
				},
				{
					Type:   lexer.TOKEN_EOF,
					Lexeme: []byte{0},
					Line:   6,
					Start:  98,
					End:    98,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := getSrcCode(t, tc.path)
			base := lexer.NewLexer(testAnnotation, src, tc.path, tc.flags)

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
