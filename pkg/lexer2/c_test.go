package lexer2_test

import (
	"fmt"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer2"
	"github.com/stretchr/testify/require"
)

func TestSingleLineComments(t *testing.T) {
	testCases := []struct {
		name             string
		testDataFilePath string
		expected         []lexer2.Token
	}{
		{
			name:             "Should return all tokens that make up the 1 single line comment in the singleline.c src code file",
			testDataFilePath: "../../testdata/singleline.c",
			expected: []lexer2.Token{
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   4,
					Start:  58,
					End:    59,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   4,
					Start:  61,
					End:    70,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("simple"),
					Line:   4,
					Start:  72,
					End:    77,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("single"),
					Line:   4,
					Start:  79,
					End:    84,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("line"),
					Line:   4,
					Start:  86,
					End:    89,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   4,
					Start:  91,
					End:    97,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   4,
					Start:  97,
					End:    97,
				},
				{
					Type:   lexer2.TOKEN_EOF,
					Lexeme: []byte{},
					Line:   6,
					Start:  123,
					End:    123,
				},
			},
		},
		{
			name:             "Should return all tokens that make up the multiple single line comments in the singleline-multi.c src code file and ignore the non annotated comment",
			testDataFilePath: "../../testdata/singleline-multi.c",
			expected: []lexer2.Token{
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   3,
					Start:  21,
					End:    22,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   3,
					Start:  24,
					End:    33,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("implement"),
					Line:   3,
					Start:  35,
					End:    43,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("the"),
					Line:   3,
					Start:  45,
					End:    47,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("main"),
					Line:   3,
					Start:  49,
					End:    52,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("function"),
					Line:   3,
					Start:  54,
					End:    61,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   3,
					Start:  61,
					End:    61,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   5,
					Start:  100,
					End:    101,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   5,
					Start:  103,
					End:    112,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("simple"),
					Line:   5,
					Start:  114,
					End:    119,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("single"),
					Line:   5,
					Start:  121,
					End:    126,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("line"),
					Line:   5,
					Start:  128,
					End:    131,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   5,
					Start:  133,
					End:    139,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   5,
					Start:  139,
					End:    139,
				},
				{
					Type:   lexer2.TOKEN_EOF,
					Lexeme: []byte{},
					Line:   10,
					Start:  246,
					End:    246,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src, fileName := getTestDataSrc(t, tc.testDataFilePath)

			base := lexer2.NewBaseLexer([]byte("@TEST_TODO"), src, fileName)
			analyzer, err := lexer2.NewLexicalAnalyzer(base)
			require.NoError(t, err)

			tokens, err := base.AnalyzeTokens(analyzer)
			require.NoError(t, err)

			require.Equal(t, tc.expected, tokens)
		})
	}
}

func TestMultiLineComments(t *testing.T) {
	testCases := []struct {
		name             string
		testDataFilePath string
		expected         []lexer2.Token
	}{
		{
			name:             "Should return all tokens that make up the multi line comments in the multiline.c src code file",
			testDataFilePath: "../../testdata/multiline.c",
			expected: []lexer2.Token{
				{
					Type:   lexer2.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   5,
					Start:  62,
					End:    63,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   5,
					Start:  65,
					End:    74,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("inline"),
					Line:   5,
					Start:  76,
					End:    81,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   5,
					Start:  83,
					End:    89,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("#1"),
					Line:   5,
					Start:  91,
					End:    92,
				},
				{
					Type:   lexer2.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte{lexer2.ASTERISK, lexer2.FORWARD_SLASH},
					Line:   5,
					Start:  94,
					End:    95,
				},
				{
					Type:   lexer2.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   6,
					Start:  115,
					End:    116,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   6,
					Start:  118,
					End:    127,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("inline"),
					Line:   6,
					Start:  129,
					End:    134,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   6,
					Start:  136,
					End:    142,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("#2"),
					Line:   6,
					Start:  144,
					End:    145,
				},
				{
					Type:   lexer2.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte{lexer2.ASTERISK, lexer2.FORWARD_SLASH},
					Line:   6,
					Start:  147,
					End:    148,
				},
				{
					Type:   lexer2.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Line:   16,
					Start:  340,
					End:    341,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   16,
					Start:  346,
					End:    355,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("drop"),
					Line:   16,
					Start:  357,
					End:    360,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("a"),
					Line:   16,
					Start:  362,
					End:    362,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("star"),
					Line:   16,
					Start:  364,
					End:    367,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("if"),
					Line:   16,
					Start:  369,
					End:    370,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("you"),
					Line:   16,
					Start:  372,
					End:    374,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("know"),
					Line:   16,
					Start:  376,
					End:    379,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("about"),
					Line:   16,
					Start:  381,
					End:    385,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("this"),
					Line:   16,
					Start:  387,
					End:    390,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("code"),
					Line:   16,
					Start:  392,
					End:    395,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("wars"),
					Line:   16,
					Start:  397,
					End:    400,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("challenge"),
					Line:   16,
					Start:  402,
					End:    410,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Digital"),
					Line:   17,
					Start:  415,
					End:    421,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Cypher"),
					Line:   17,
					Start:  423,
					End:    428,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("assigns"),
					Line:   17,
					Start:  430,
					End:    436,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("to"),
					Line:   17,
					Start:  438,
					End:    439,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("each"),
					Line:   17,
					Start:  441,
					End:    444,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("letter"),
					Line:   17,
					Start:  446,
					End:    451,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("of"),
					Line:   17,
					Start:  453,
					End:    454,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   17,
					Start:  456,
					End:    458,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("alphabet"),
					Line:   17,
					Start:  460,
					End:    467,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("unique"),
					Line:   17,
					Start:  469,
					End:    474,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("number."),
					Line:   17,
					Start:  476,
					End:    482,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Instead"),
					Line:   18,
					Start:  487,
					End:    493,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("of"),
					Line:   18,
					Start:  495,
					End:    496,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("letters"),
					Line:   18,
					Start:  498,
					End:    504,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("in"),
					Line:   18,
					Start:  506,
					End:    507,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("encrypted"),
					Line:   18,
					Start:  509,
					End:    517,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("word"),
					Line:   18,
					Start:  519,
					End:    522,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("we"),
					Line:   18,
					Start:  524,
					End:    525,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("write"),
					Line:   18,
					Start:  527,
					End:    531,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   18,
					Start:  533,
					End:    535,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("corresponding"),
					Line:   18,
					Start:  537,
					End:    549,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("number"),
					Line:   18,
					Start:  551,
					End:    556,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("Then"),
					Line:   19,
					Start:  561,
					End:    564,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("we"),
					Line:   19,
					Start:  566,
					End:    567,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("add"),
					Line:   19,
					Start:  569,
					End:    571,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("to"),
					Line:   19,
					Start:  573,
					End:    574,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("each"),
					Line:   19,
					Start:  576,
					End:    579,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("obtained"),
					Line:   19,
					Start:  581,
					End:    588,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("digit"),
					Line:   19,
					Start:  590,
					End:    594,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("consecutive"),
					Line:   19,
					Start:  596,
					End:    606,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("digits"),
					Line:   19,
					Start:  608,
					End:    613,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("from"),
					Line:   19,
					Start:  615,
					End:    618,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("the"),
					Line:   19,
					Start:  620,
					End:    622,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("key"),
					Line:   19,
					Start:  624,
					End:    626,
				},
				{
					Type:   lexer2.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte{lexer2.ASTERISK, lexer2.FORWARD_SLASH},
					Line:   21,
					Start:  631,
					End:    632,
				},
				{
					Type:   lexer2.TOKEN_EOF,
					Lexeme: []byte{},
					Line:   34,
					Start:  977,
					End:    977,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src, fileName := getTestDataSrc(t, tc.testDataFilePath)

			base := lexer2.NewBaseLexer([]byte("@TEST_TODO"), src, fileName)
			analyzer, err := lexer2.NewLexicalAnalyzer(base)
			require.NoError(t, err)

			tokens, err := base.AnalyzeTokens(analyzer)
			require.NoError(t, err)

			for _, token := range tokens {
				if token.Line == 6 {
					fmt.Println(string(token.Lexeme))
				}
			}

			require.Equal(t, tc.expected, tokens)
		})
	}
}
