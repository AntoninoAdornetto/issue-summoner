package lexer3_test

import (
	"fmt"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer3"
	"github.com/stretchr/testify/require"
)

func TestSingleLineComments(t *testing.T) {
	testCases := []struct {
		name             string
		testDataFilePath string
		expected         []lexer3.Token
	}{
		{
			name:             "Should return all tokens that make up the 1 single line comment in the singleline.c src code file",
			testDataFilePath: "../../testdata/singleline.c",
			expected: []lexer3.Token{
				{
					Type:   lexer3.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   4,
					Start:  58,
					End:    59,
				},
				{
					Type:   lexer3.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   4,
					Start:  61,
					End:    70,
				},
				{
					Type:   lexer3.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("simple"),
					Line:   4,
					Start:  72,
					End:    77,
				},
				{
					Type:   lexer3.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("single"),
					Line:   4,
					Start:  79,
					End:    84,
				},
				{
					Type:   lexer3.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("line"),
					Line:   4,
					Start:  86,
					End:    89,
				},
				{
					Type:   lexer3.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   4,
					Start:  91,
					End:    97,
				},
				{
					Type:   lexer3.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{lexer3.NEWLINE},
					Line:   4,
					Start:  98,
					End:    98,
				},
				{
					Type:   lexer3.TOKEN_EOF,
					Lexeme: []byte{},
					Line:   5,
					Start:  123,
					End:    123,
				},
			},
		},
		// {
		// 	name:             "Should return all tokens that make up the multiple single line comments in the singleline-multi.c src code file and ignore the non annotated comment",
		// 	testDataFilePath: "../../testdata/singleline-multi.c",
		// 	expected: []lexer3.Token{
		// 		{
		// 			Type:   lexer3.TOKEN_SINGLE_LINE_COMMENT_START,
		// 			Lexeme: []byte("//"),
		// 			Line:   3,
		// 			Start:  21,
		// 			End:    22,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_ANNOTATION,
		// 			Lexeme: []byte("@TEST_TODO"),
		// 			Line:   3,
		// 			Start:  24,
		// 			End:    33,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("implement"),
		// 			Line:   3,
		// 			Start:  35,
		// 			End:    43,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("the"),
		// 			Line:   3,
		// 			Start:  45,
		// 			End:    47,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("main"),
		// 			Line:   3,
		// 			Start:  49,
		// 			End:    52,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("function"),
		// 			Line:   3,
		// 			Start:  54,
		// 			End:    61,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_SINGLE_LINE_COMMENT_END,
		// 			Lexeme: []byte{0},
		// 			Line:   3,
		// 			Start:  61,
		// 			End:    61,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_SINGLE_LINE_COMMENT_START,
		// 			Lexeme: []byte("//"),
		// 			Line:   5,
		// 			Start:  100,
		// 			End:    101,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_ANNOTATION,
		// 			Lexeme: []byte("@TEST_TODO"),
		// 			Line:   5,
		// 			Start:  103,
		// 			End:    112,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("simple"),
		// 			Line:   5,
		// 			Start:  114,
		// 			End:    119,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("single"),
		// 			Line:   5,
		// 			Start:  121,
		// 			End:    126,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("line"),
		// 			Line:   5,
		// 			Start:  128,
		// 			End:    131,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_COMMENT_TITLE,
		// 			Lexeme: []byte("comment"),
		// 			Line:   5,
		// 			Start:  133,
		// 			End:    139,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_SINGLE_LINE_COMMENT_END,
		// 			Lexeme: []byte{0},
		// 			Line:   5,
		// 			Start:  139,
		// 			End:    139,
		// 		},
		// 		{
		// 			Type:   lexer3.TOKEN_EOF,
		// 			Lexeme: []byte{},
		// 			Line:   10,
		// 			Start:  246,
		// 			End:    246,
		// 		},
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src, fileName := getTestDataSrc(t, tc.testDataFilePath)

			base := lexer3.NewBaseLexer([]byte("@TEST_TODO"), src, fileName)
			analyzer, err := lexer3.NewLexicalAnalyzer(base)
			require.NoError(t, err)

			tokens, err := base.AnalyzeTokens(analyzer)
			require.NoError(t, err)

			for _, token := range tokens {
				fmt.Println(string(token.Lexeme))
			}

			require.Equal(t, tc.expected, tokens)
		})
	}
}
