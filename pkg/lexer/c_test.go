package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeTokenSingleLineComments(t *testing.T) {
	testCases := []struct {
		name     string
		srcCode  []byte
		fileName string
		expected []lexer.Token
	}{
		{
			name: "should not create any tokens when consuming a non-comment notation byte",
			// int is not "//" or "/*" - denotes the opening notation of a comment in c-like languages
			srcCode:  []byte("int"),
			fileName: "main.c",
			expected: []lexer.Token{},
		},
		{
			name:     "should not create any tokens when consuming single line comment bytes that do not have an issue annotation",
			srcCode:  []byte("// regular single line comment with no issue annotation in c"),
			fileName: "main.c",
			expected: []lexer.Token{},
		},
		{
			name:     "should create the comment start, comment annotation and comment end tokens",
			srcCode:  []byte("// @TEST_ANNOTATION"),
			fileName: "main.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  3,
					End:    18,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Start:  19,
					End:    19,
					Line:   1,
				},
			},
		},
		{
			name:     "should create the comment start, comment annotation, comment title and comment end tokens",
			srcCode:  []byte("// @TEST_ANNOTATION check for edge cases"),
			fileName: "main.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  3,
					End:    18,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("check"),
					Start:  20,
					End:    24,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("for"),
					Start:  26,
					End:    28,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("edge"),
					Start:  30,
					End:    33,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("cases"),
					Start:  35,
					End:    39,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Start:  40,
					End:    40,
					Line:   1,
				},
			},
		},
	}

	for _, tc := range testCases {
		baseLexer := lexer.NewLexer(testAnnotation, tc.srcCode, tc.fileName)
		cLexer := lexer.Clexer{Base: baseLexer}
		err := cLexer.AnalyzeToken()
		require.NoError(t, err)
		require.Equal(t, tc.expected, baseLexer.Tokens)
	}
}
